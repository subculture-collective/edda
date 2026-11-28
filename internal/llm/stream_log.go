package llm

import (
	"context"
	"errors"
	"log/slog"

	"git.subcult.tv/subculture-collective/edda/internal/logging"
)

// logStreamFailure records why a stream closed without Done. The channel
// carries no error, so this record is the only trace of the cause. A caller
// that canceled the request gave up on the stream, which is not a failure.
func logStreamFailure(ctx context.Context, logger *slog.Logger, model string, err error) {
	if errors.Is(ctx.Err(), context.Canceled) {
		return
	}
	logger.Error("stream failed", "model", model, "error", logging.Redact(err.Error()))
}
