// Command logs tails a edda JSONL log file with colorized output and filtering.
package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"git.subcult.tv/subculture-collective/edda/internal/logging"
)

func main() {
	filePath := flag.String("file", ".logs/edda.jsonl", "path to JSONL log file")
	levelStr := flag.String("level", "debug", "minimum log level: debug, info, warn, error")
	var serviceArgs multiValueFlag
	flag.Var(&serviceArgs, "service", "service filter; repeat or comma-separate values (empty = all)")
	history := flag.Int("history", 0, "number of historical lines to show before tailing (0 = skip to end)")
	flag.Parse()

	minLevel := parseLevel(*levelStr)
	services := parseServices(serviceArgs)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	file := waitForFile(ctx, *filePath)
	if file == nil {
		return // context cancelled while waiting
	}
	defer file.Close()

	if *history > 0 {
		printHistory(file, *history, minLevel, services)
	} else {
		// Seek to end for live tail only.
		if _, err := file.Seek(0, io.SeekEnd); err != nil {
			fmt.Fprintf(os.Stderr, "seek: %v\n", err)
			os.Exit(1)
		}
	}

	tail(ctx, file, os.Stdout, minLevel, services)
}

// waitForFile opens the file, polling every 500ms if it doesn't exist yet.
func waitForFile(ctx context.Context, path string) *os.File {
	f, err := os.Open(path)
	if err == nil {
		return f
	}
	if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "open %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Waiting for log file...\n")
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			f, err = os.Open(path)
			if err == nil {
				return f
			}
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "open %s: %v\n", path, err)
				os.Exit(1)
			}
		}
	}
}

// printHistory seeks backwards from end-of-file to find the last N lines, then prints them.
func printHistory(file *os.File, n int, minLevel slog.Level, services map[string]bool) {
	// Read entire file to find last N lines. For simplicity and correctness
	// with variable-width UTF-8 JSONL, we scan forward and keep a ring of lines.
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		fmt.Fprintf(os.Stderr, "seek: %v\n", err)
		return
	}

	// Collect all lines, then take last N that pass filters.
	var matched []logging.LogEntry
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		entry, err := logging.ParseJSONLEntry(line)
		if err != nil {
			continue
		}
		if !passesFilter(entry, minLevel, services) {
			continue
		}
		matched = append(matched, entry)
	}

	// Print last N.
	start := 0
	if len(matched) > n {
		start = len(matched) - n
	}
	for _, entry := range matched[start:] {
		fmt.Print(logging.FormatColorLine(entry))
	}

	// Leave file position at end for subsequent tailing.
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		fmt.Fprintf(os.Stderr, "seek: %v\n", err)
	}
}

// tail polls the file for new JSONL lines and prints matching entries to out.
// A bufio.Scanner stops for good at its first EOF, so tail reads with a
// bufio.Reader, keeps a partial last line until its newline arrives and polls
// at EOF. When the file shrinks below the read offset it was truncated or
// replaced in place, and tail starts again from its beginning. When the path
// names a different file, the log was rotated by rename: tail reads the old
// file to its end, then reads the new one from its beginning. tail takes
// ownership of file and closes every file it holds.
func tail(ctx context.Context, file *os.File, out io.Writer, minLevel slog.Level, services map[string]bool) {
	path := file.Name()
	defer func() { _ = file.Close() }()
	offset, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seek: %v\n", err)
		return
	}
	reader := bufio.NewReader(file)
	var partial []byte
	var next *os.File // the file now at path, while the old one is drained
	defer func() {
		if next != nil {
			_ = next.Close()
		}
	}()

	for {
		chunk, err := reader.ReadBytes('\n')
		offset += int64(len(chunk))
		partial = append(partial, chunk...)
		if err == nil {
			printLine(out, partial, minLevel, services)
			partial = partial[:0]
			continue
		}
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "read: %v\n", err)
			return
		}

		if next != nil {
			// The old file is drained. No more of its last line will come.
			printLine(out, partial, minLevel, services)
			partial = partial[:0]
			_ = file.Close()
			file, next, offset = next, nil, 0
			reader.Reset(file)
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
		}

		// Read once more before leaving the old file: lines may have been
		// added to it between the EOF above and the rename.
		if next = openIfReplaced(file, path); next != nil {
			continue
		}

		if info, err := file.Stat(); err == nil && info.Size() < offset {
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				fmt.Fprintf(os.Stderr, "seek: %v\n", err)
				return
			}
			offset = 0
			partial = partial[:0]
			reader.Reset(file)
		}
	}
}

// openIfReplaced opens the file at path when it is not the open file. It
// returns nil while path names the open file, or nothing that can be opened,
// which is normal between a rename and the creation of the new file.
func openIfReplaced(file *os.File, path string) *os.File {
	atPath, err := os.Stat(path)
	if err != nil {
		return nil
	}
	open, err := file.Stat()
	if err != nil || os.SameFile(open, atPath) {
		return nil
	}
	replacement, err := os.Open(path)
	if err != nil {
		return nil
	}
	return replacement
}

func printLine(out io.Writer, line []byte, minLevel slog.Level, services map[string]bool) {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return
	}
	entry, err := logging.ParseJSONLEntry(line)
	if err != nil || !passesFilter(entry, minLevel, services) {
		return
	}
	fmt.Fprint(out, logging.FormatColorLine(entry))
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

type multiValueFlag []string

func (f *multiValueFlag) String() string {
	if f == nil {
		return ""
	}
	return strings.Join(*f, ",")
}

func (f *multiValueFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func parseServices(values []string) map[string]bool {
	if len(values) == 0 {
		return nil
	}
	m := make(map[string]bool)
	for _, value := range values {
		for _, service := range strings.Split(value, ",") {
			service = strings.ToLower(strings.TrimSpace(service))
			if service == "" {
				continue
			}
			m[service] = true
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func passesFilter(entry logging.LogEntry, minLevel slog.Level, services map[string]bool) bool {
	if entry.Level < minLevel {
		return false
	}
	if len(services) > 0 && !services[strings.ToLower(entry.Service)] {
		return false
	}
	return true
}
