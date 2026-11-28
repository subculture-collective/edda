package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackc/pgx/v5/pgxpool"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
	"git.subcult.tv/subculture-collective/edda/internal/logging"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/tui"
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

	charm := logging.NewCharmLogger(logResult.BridgeWriter)
	charm.SetDefault() // without a prefix: each package adds its own
	logger := charm.WithPrefix("tui")

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Error("load config", "error", err)
		return 1
	}
	if _, err := llm.NewLLMProvider(cfg); err != nil {
		logger.Error("initialize llm provider", "error", err)
		return 1
	}

	logStartup(logger, cfg)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		logger.Error("open database", "error", err)
		return 1
	}
	defer pool.Close()

	queries := statedb.New(pool)
	provider, err := newLLMProvider(cfg)
	if err != nil {
		logger.Error("create llm provider", "error", err)
		return 1
	}
	gameEngine, err := engine.New(pool, provider, cfg.LLM, engine.WithLogger(slog.Default().WithGroup("engine")))
	if err != nil {
		logger.Error("create engine", "error", err)
		return 1
	}

	programOptions := []tea.ProgramOption{
		tea.WithContext(ctx),
	}
	if tuiAltScreenEnabled() {
		programOptions = append(programOptions, tea.WithAltScreen())
	}
	if tuiMouseEnabled() {
		programOptions = append(programOptions, tea.WithMouseCellMotion())
	}

	p := tea.NewProgram(
		tui.NewLauncherWithEngine(
			cfg,
			ctx,
			queries,
			gameEngine,
			tui.WithLLMProvider(provider),
			tui.WithLogBuffer(logResult.RingBuffer),
		),
		programOptions...,
	)

	go func() {
		<-ctx.Done()
		logger.Info("shutdown signal received")
	}()

	if err := runQuietly(logResult, func() error {
		_, err := p.Run()
		return err
	}); err != nil {
		if ctx.Err() != nil && (errors.Is(err, tea.ErrInterrupted) || errors.Is(err, tea.ErrProgramKilled)) {
			logger.Info("TUI shutdown complete")
			return 0
		}
		logger.Error("tui error", "error", err)
		return 1
	}

	logger.Info("TUI stopped")
	return 0
}

// logStartup logs the active provider's settings, with any credentials removed
// from its endpoint.
func logStartup(logger logging.Logger, cfg config.Config) {
	keyvals := []any{"provider", cfg.LLM.Provider}
	switch cfg.LLM.Provider {
	case "ollama":
		keyvals = append(keyvals,
			"model", cfg.LLM.Ollama.Model,
			"endpoint", logging.RedactURL(cfg.LLM.Ollama.Endpoint),
			"timeout", cfg.LLM.Ollama.RequestTimeout(),
		)
	case "claude":
		keyvals = append(keyvals, "model", cfg.LLM.Claude.Model)
	case "openrouter":
		keyvals = append(keyvals, "model", cfg.LLM.OpenRouter.Model)
	case "chatgpt":
		keyvals = append(keyvals, "model", cfg.LLM.ChatGPT.Model)
	case "claudecode":
		keyvals = append(keyvals, "endpoint", logging.RedactURL(cfg.LLM.ClaudeCode.Endpoint))
	}
	keyvals = append(keyvals, "altscreen", tuiAltScreenEnabled(), "mouse", tuiMouseEnabled())
	logger.Info("starting TUI", keyvals...)
}

// runQuietly runs the terminal UI with the stderr copy of log records off. The
// log panel and the log file still receive them.
func runQuietly(logs *logging.SetupResult, run func() error) error {
	logs.SetConsole(false)
	defer logs.SetConsole(true)
	return run()
}

func parseConfigPath(args []string, defaultPath string) (string, error) {
	fs := flag.NewFlagSet("tui", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	configPath := defaultPath
	fs.StringVar(&configPath, "config", configPath, "Path to config file")

	if err := fs.Parse(args); err != nil {
		return "", err
	}
	return configPath, nil
}

func newLLMProvider(cfg config.Config) (llm.Provider, error) {
	switch cfg.LLM.Provider {
	case "claude":
		return llm.NewClaudeClient("", cfg.LLM.Claude.APIKey, cfg.LLM.Claude.Model), nil
	case "ollama":
		return llm.NewOllamaClientWithTimeout(cfg.LLM.Ollama.Endpoint, cfg.LLM.Ollama.Model, cfg.LLM.Ollama.RequestTimeout()).
			WithAPIKey(cfg.LLM.Ollama.APIKey), nil
	default:
		return nil, fmt.Errorf("unknown llm provider: %q", cfg.LLM.Provider)
	}
}

func tuiAltScreenEnabled() bool {
	return os.Getenv("EDDA_TUI_ALTSCREEN") != "0"
}

func tuiMouseEnabled() bool {
	return os.Getenv("EDDA_TUI_MOUSE") != "0"
}
