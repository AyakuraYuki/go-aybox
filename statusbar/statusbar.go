package statusbar

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/term"
	"golang.org/x/text/width"
)

// Style controls the visual presentation of the status bar.
type Style int

const (
	StyleASCII Style = iota // plain ASCII separators
	StyleEmoji              // emoji icons as field indicators
)

// StatusBar renders a single-line status display in the terminal showing the
// current time, elapsed time, and an optional message. It runs in a background
// goroutine and serializes all terminal writes through a mutex so that log
// output interleaved via Writer() does not corrupt the display.
//
// Typical usage with the log module:
//
//	bar := statusbar.New()
//	cw  := console.New(console.WithWriter(bar.Writer()))
//	logger := log.New(log.WithWriters(cw))
//	bar.Start()
//	defer bar.Stop()
type StatusBar struct {
	mu            sync.Mutex
	startTime     time.Time
	message       string
	out           io.Writer
	interval      time.Duration
	stopCh        chan struct{}
	running       bool
	style         Style
	fgColor       int
	bgColor       int
	width         int  // 0: no constraint; overridden by autoWidth at render time
	autoWidth     bool // query terminal width on each render
	fixedToBottom bool // always render on the last terminal line
}

// New creates a StatusBar. Call Start to begin rendering.
// Defaults: StyleEmoji, BgCyan background, FgBlack foreground, 1s refresh,
// auto-width.
func New(opts ...Option) *StatusBar {
	s := &StatusBar{
		startTime: time.Now(),
		out:       os.Stdout,
		interval:  time.Second,
		stopCh:    make(chan struct{}),
		style:     StyleEmoji,
		fgColor:   FgBlack,
		bgColor:   BgCyan,
		autoWidth: true,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Start begins rendering the status bar in a background goroutine.
// Calling Start on an already-running bar is a no-op.
func (s *StatusBar) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()
	go s.loop()
}

// Stop halts the background goroutine, renders the final frame, and moves the
// cursor to a new line so the bar remains visible and subsequent output starts
// cleanly. Calling Stop on a bar that is not running is a no-op.
func (s *StatusBar) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopCh)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.fixedToBottom {
		// Position at the last line so the final frame lands there.
		if h := queryTerminalHeight(); h > 0 {
			_, _ = fmt.Fprintf(s.out, "\033[%d;1H", h)
		}
		content := s.buildContent()
		_, _ = fmt.Fprintf(s.out, "\033[%dm\033[%dm%s\033[K\033[0m\n",
			s.bgColor, s.fgColor, content)
	} else {
		// renderLocked writes \r first, overwriting whatever the ticker last
		// rendered, so the final frame is never duplicated.
		s.renderLocked()
		_, _ = fmt.Fprint(s.out, "\n")
	}
}

// SetMessage sets the message shown on the status line.
// An empty string clears the message.
func (s *StatusBar) SetMessage(msg string) {
	s.mu.Lock()
	s.message = msg
	s.mu.Unlock()
}

// ClearMessage removes the message from the status line.
func (s *StatusBar) ClearMessage() {
	s.mu.Lock()
	s.message = ""
	s.mu.Unlock()
}

// loop is the background goroutine that periodically redraws the status line.
func (s *StatusBar) loop() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.mu.Lock()
	s.renderLocked()
	s.mu.Unlock()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()
			s.renderLocked()
			s.mu.Unlock()
		}
	}
}

// buildContent builds the formatted status line text without ANSI color codes.
// Must be called with s.mu held.
func (s *StatusBar) buildContent() string {
	now := time.Now().Format("15:04:05")
	elapsed := formatElapsed(time.Since(s.startTime))

	var content string
	switch s.style {
	case StyleEmoji:
		content = fmt.Sprintf(" 🕐 %s | ⏱ %s", now, elapsed)
		if s.message != "" {
			content += fmt.Sprintf(" | 💬 %s ", s.message)
		} else {
			content += " "
		}
	default: // StyleASCII
		content = fmt.Sprintf(" %s | %s", now, elapsed)
		if s.message != "" {
			content += fmt.Sprintf(" | %s ", s.message)
		} else {
			content += " "
		}
	}

	effectiveWidth := s.width
	if s.autoWidth {
		if w := queryTerminalWidth(); w > 0 {
			effectiveWidth = w
		}
	}
	if effectiveWidth > 0 {
		content = fitWidth(content, effectiveWidth)
	}
	return content
}

// renderLocked writes the status line. Must be called with s.mu held.
// In fixedToBottom mode the cursor is saved, the bar is drawn on the last
// terminal line, and the cursor is restored so ongoing output is unaffected.
func (s *StatusBar) renderLocked() {
	bar := fmt.Sprintf("\033[%dm\033[%dm%s\033[K\033[0m",
		s.bgColor, s.fgColor, s.buildContent())

	if s.fixedToBottom {
		if h := queryTerminalHeight(); h > 0 {
			// Save cursor → move to last line → render → restore cursor.
			_, _ = fmt.Fprintf(s.out, "\033[s\033[%d;1H%s\033[u", h, bar)
			return
		}
		// Fall through to normal render when terminal height is unavailable.
	}

	_, _ = fmt.Fprintf(s.out, "\r%s", bar)
}

// formatElapsed formats d as mm:ss, or hh:mm:ss when hours are non-zero.
func formatElapsed(d time.Duration) string {
	d = d.Truncate(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	sec := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, sec)
	}
	return fmt.Sprintf("%02d:%02d", m, sec)
}

// queryTerminalHeight probes the terminal height via term.GetSize on stdout,
// stderr, and stdin in order, then falls back to the LINES environment variable.
// Returns 0 when all probes fail.
func queryTerminalHeight() int {
	for _, f := range []*os.File{os.Stdout, os.Stderr, os.Stdin} {
		if _, h, err := term.GetSize(int(f.Fd())); err == nil && h > 0 {
			return h
		}
	}
	if lines := os.Getenv("LINES"); lines != "" {
		if h, err := strconv.Atoi(lines); err == nil && h > 0 {
			return h
		}
	}
	return 0
}

// queryTerminalWidth probes the terminal width via term.GetSize on stdout,
// stderr, and stdin in order (covers native shells on Windows/macOS/Linux and
// JetBrains IDE run windows backed by a PTY), then falls back to the COLUMNS
// environment variable. Returns 0 when all probes fail.
func queryTerminalWidth() int {
	for _, f := range []*os.File{os.Stdout, os.Stderr, os.Stdin} {
		if w, _, err := term.GetSize(int(f.Fd())); err == nil && w > 0 {
			return w
		}
	}
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if w, err := strconv.Atoi(cols); err == nil && w > 0 {
			return w
		}
	}
	return 0
}

// runeDisplayWidth returns the terminal column count for r.
// Delegates to golang.org/x/text/width (Unicode East Asian Width):
// Wide/Fullwidth (emoji, CJK) → 2; control chars → 0; others → 1.
func runeDisplayWidth(r rune) int {
	if r < 0x20 || r == 0x7F {
		return 0
	}
	switch width.LookupRune(r).Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	default:
		return 1
	}
}

// fitWidth truncates or right-pads s to exactly n terminal columns.
// When a wide character would overshoot the boundary, spaces fill the gap.
func fitWidth(s string, n int) string {
	cur := 0
	for i, r := range s {
		rw := runeDisplayWidth(r)
		if cur+rw > n {
			result := s[:i]
			for cur < n {
				result += " "
				cur++
			}
			return result
		}
		cur += rw
	}
	for cur < n {
		s += " "
		cur++
	}
	return s
}

// Writer returns an io.Writer that serializes writes with the status bar's
// mutex. Before each write the current status line is cleared; after the write
// the status bar is redrawn. Pass this as the underlying writer of the log
// console writer to prevent display corruption:
//
//	cw := console.New(console.WithWriter(bar.Writer()))
func (s *StatusBar) Writer() io.Writer {
	return &statusWriter{bar: s}
}

// statusWriter clears and redraws the status bar around every Write call,
// ensuring log lines and the status bar do not overwrite each other.
type statusWriter struct {
	bar *StatusBar
}

func (w *statusWriter) Write(p []byte) (n int, err error) {
	w.bar.mu.Lock()
	defer w.bar.mu.Unlock()

	// Clear the current status line, write the log content, then redraw.
	// In fixedToBottom mode renderLocked handles cursor save/restore so the
	// bar is redrawn at the last line without disturbing the log output.
	_, _ = fmt.Fprint(w.bar.out, "\r\033[K")
	n, err = w.bar.out.Write(p)
	if err != nil {
		return
	}
	if w.bar.running {
		w.bar.renderLocked()
	}
	return
}
