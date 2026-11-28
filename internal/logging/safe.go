package logging

import (
	"fmt"

	"github.com/charmbracelet/log"
)

// Logger is what packages need in order to log through a command's logger.
type Logger interface {
	Error(msg any, keyvals ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

// CharmLogger is a charmbracelet/log logger that never panics on a logged value.
// The JSON formatter calls Error, String, LogValue and MarshalJSON without
// recovering, so a typed nil error would crash the caller. CharmLogger turns keys
// and values into strings with fmt first, as the text formatter does: fmt
// recovers those panics. The wrapped logger is private so that a method cannot be
// reached without this step.
type CharmLogger struct {
	base *log.Logger
}

// WrapCharmLogger makes logging through base safe for any value.
func WrapCharmLogger(base *log.Logger) *CharmLogger {
	return &CharmLogger{base: base}
}

// SetDefault makes the wrapped logger charmbracelet/log's default. Package-level
// calls such as log.Error reach it directly, so give them no key/value pairs;
// pass log.Default() to WrapCharmLogger for those.
func (l *CharmLogger) SetDefault() { log.SetDefault(l.base) }

// With returns a logger that adds keyvals to every record.
func (l *CharmLogger) With(keyvals ...any) *CharmLogger {
	return WrapCharmLogger(l.base.With(safeKeyvals(keyvals)...))
}

// WithPrefix returns a logger with the given prefix.
func (l *CharmLogger) WithPrefix(prefix string) *CharmLogger {
	return WrapCharmLogger(l.base.WithPrefix(prefix))
}

// Log writes a record at level. charmbracelet/log formats msg with fmt.
func (l *CharmLogger) Log(level log.Level, msg any, keyvals ...any) {
	l.base.Log(level, msg, safeKeyvals(keyvals)...)
}

// Debug writes a debug record.
func (l *CharmLogger) Debug(msg any, keyvals ...any) { l.Log(log.DebugLevel, msg, keyvals...) }

// Info writes an info record.
func (l *CharmLogger) Info(msg any, keyvals ...any) { l.Log(log.InfoLevel, msg, keyvals...) }

// Warn writes a warning record.
func (l *CharmLogger) Warn(msg any, keyvals ...any) { l.Log(log.WarnLevel, msg, keyvals...) }

// Error writes an error record.
func (l *CharmLogger) Error(msg any, keyvals ...any) { l.Log(log.ErrorLevel, msg, keyvals...) }

// Debugf writes a formatted debug record.
func (l *CharmLogger) Debugf(format string, args ...any) { l.base.Debugf(format, args...) }

// Infof writes a formatted info record.
func (l *CharmLogger) Infof(format string, args ...any) { l.base.Infof(format, args...) }

// Errorf writes a formatted error record.
func (l *CharmLogger) Errorf(format string, args ...any) { l.base.Errorf(format, args...) }

// safeKeyvals returns a copy of keyvals with each key and value written as the
// text formatter writes it.
func safeKeyvals(keyvals []any) []any {
	safe := make([]any, len(keyvals))
	for i, v := range keyvals {
		if i%2 == 0 {
			safe[i] = fmt.Sprint(v)
		} else {
			safe[i] = fmt.Sprintf("%+v", v)
		}
	}
	return safe
}
