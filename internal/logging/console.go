package logging

import (
	"io"
	"sync/atomic"
)

// consoleWriter forwards to w unless muted.
type consoleWriter struct {
	w     io.Writer
	muted atomic.Bool
}

func (c *consoleWriter) Write(p []byte) (int, error) {
	if c.muted.Load() {
		return len(p), nil
	}
	return c.w.Write(p)
}
