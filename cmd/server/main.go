package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"

	"git.subcult.tv/subculture-collective/edda/internal/assembly"
	"git.subcult.tv/subculture-collective/edda/internal/bootstrap"
	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
	"git.subcult.tv/subculture-collective/edda/internal/logging"
	"git.subcult.tv/subculture-collective/edda/internal/memory"
	"git.subcult.tv/subculture-collective/edda/internal/saves"
	serverroutes "git.subcult.tv/subculture-collective/edda/internal/server"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	configPath, err := parseConfigPath(args, os.Getenv("EDDA_CONFIG"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse flags: %v\n", err)
		return 2
	}

	logResult, err := logging.Setup(".logs/edda.jsonl", slog.LevelDebug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init logging: %v\n", err)
		return 1
	}
	defer logResult.Cleanup()

	logger := log.NewWithOptions(logResult.BridgeWriter, log.Options{ReportTimestamp: true})
	log.SetDefault(logger)

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Errorf("load config: %v", err)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		logger.Errorf("open database: %v", err)
		return 1
	}
	defer pool.Close()

	queries := statedb.New(pool)

	bootResult, err := bootstrap.Run(ctx, queries)
	if err != nil {
		logger.Errorf("bootstrap: %v", err)
		return 1
	}
	defaultUserID := bootResult.User.ID.Bytes

	roster, err := llm.NewRoster(cfg)
	if err != nil {
		logger.Errorf("create llm roster: %v", err)
		return 1
	}
	provider := roster.Provider(llm.FlowDefault)
	gmProvider := roster.Provider(llm.FlowGMTurn)

	saveStore := saves.NewStore(pool)

	engineOpts := []engine.Option{
		engine.WithLogger(slog.Default().WithGroup("engine")),
		engine.WithSaveStore(saveStore),
		engine.WithLLMProviders(engine.EngineLLMProviders{
			PostTurnState:      roster.Provider(llm.FlowPostTurnState),
			ChoiceFallback:     roster.Provider(llm.FlowChoiceFallback),
			ContextTokenBudget: roster.Info(llm.FlowGMTurn).ContextTokenBudget,
		}),
	}
	embedding := cfg.LLM.MemoryEmbedding
	if !embedding.Enabled && cfg.LLM.Provider == "ollama" && cfg.LLM.Ollama.EmbeddingModel != "" {
		embedding = config.MemoryEmbeddingConfig{Enabled: true, Endpoint: cfg.LLM.Ollama.EmbeddingEndpoint, Model: cfg.LLM.Ollama.EmbeddingModel, Dimension: cfg.LLM.Ollama.EmbeddingDimension, APIKey: cfg.LLM.Ollama.APIKey}
	}
	if embedding.Enabled {
		if embedding.Dimension <= 0 {
			embedding.Dimension = memory.DefaultVectorDimension
		}
		if err := memory.ValidateMemoryEmbeddingDimension(ctx, pool, embedding.Dimension); err != nil {
			logger.Errorf("validate memory embedding dimension: %v", err)
			return 1
		}
		embedEndpoint := embedding.Endpoint
		if embedEndpoint == "" {
			embedEndpoint = cfg.LLM.Ollama.Endpoint
		}
		embedder := memory.NewOllamaEmbedder(
			embedEndpoint, embedding.Model,
			memory.WithOllamaEmbedderDimension(embedding.Dimension),
			memory.WithOllamaEmbedderTimeout(cfg.LLM.Ollama.RequestTimeout()),
			memory.WithOllamaEmbedderAPIKey(embedding.APIKey),
		)
		searcher := memory.NewSearcher(embedder, queries)
		tier3 := assembly.NewTier3Retriever(searcher, 5, slog.Default().WithGroup("tier3"))
		engineOpts = append(engineOpts,
			engine.WithTier3Retriever(tier3),
			engine.WithEmbedder(embedder),
			engine.WithSearcher(searcher),
		)
	} else {
		logger.Warn("semantic campaign memory unavailable: llm.memoryembedding.enabled is false")
	}

	gameEngine, err := engine.New(pool, gmProvider, cfg.LLM, engineOpts...)
	if err != nil {
		logger.Errorf("create engine: %v", err)
		return 1
	}
	router := serverroutes.NewRouterWithProvider(logger, gameEngine, queries, provider, pool, defaultUserID, cfg, saveStore)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	logger.Infof("starting HTTP server on %s (provider=%s)", addr, cfg.LLM.Provider)

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("server failed: %v", err)
			return 1
		}
	case <-ctx.Done():
		logger.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("graceful shutdown failed: %v", err)
			return 1
		}

		if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("server failed during shutdown: %v", err)
			return 1
		}
	}

	logger.Info("server shutdown complete")
	return 0
}

func parseConfigPath(args []string, defaultPath string) (string, error) {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	configPath := defaultPath
	fs.StringVar(&configPath, "config", configPath, "Path to config file")

	if err := fs.Parse(args); err != nil {
		return "", err
	}
	return configPath, nil
}
