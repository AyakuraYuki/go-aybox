package statusbar

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/term"
)

// Style controls the visual presentation of the status bar.
type Style int

const (
	// StyleASCII uses plain ASCII separators.
	StyleASCII Style = iota
	// StyleEmoji uses emoji icons as field indicators.
	StyleEmoji
)

// ANSI foreground color codes.
const (
	FgBlack   = 30
	FgRed     = 31
	FgGreen   = 32
	FgYellow  = 33
	FgBlue    = 34
	FgMagenta = 35
	FgCyan    = 36
	FgWhite   = 37
)

// ANSI background color codes.
const (
	BgBlack   = 40
	BgRed     = 41
	BgGreen   = 42
	BgYellow  = 43
	BgBlue    = 44
	BgMagenta = 45
	BgCyan    = 46
	BgWhite   = 47
)

// StatusBar renders a single-line status display in the terminal showing the
// current time, program elapsed time, and an optional message. It runs in a
// background goroutine and serializes all terminal writes through a mutex so
// that log output interleaved via Writer() does not corrupt the display.
//
// Typical usage with the log module:
//
//	bar := statusbar.New()
//	cw  := console.New(console.WithWriter(bar.Writer()))
//	logger := log.New(log.WithWriters(cw))
//	bar.Start()
//	defer bar.Stop()
type StatusBar struct {
	mu        sync.Mutex
	startTime time.Time
	message   string
	out       io.Writer
	interval  time.Duration
	stopCh    chan struct{}
	running   bool
	style     Style
	fgColor   int
	bgColor   int
	width     int  // 0 means no constraint; overridden by autoWidth at render time
	autoWidth bool // when true, query terminal width on each render
}

// Option configures a StatusBar at construction time.
type Option func(*StatusBar)

// WithOutput sets the output writer (default: os.Stdout).
func WithOutput(w io.Writer) Option {
	return func(s *StatusBar) {
		if w != nil {
			s.out = w
		}
	}
}

// WithRefreshRate sets how often the status line is redrawn (default: 1s).
func WithRefreshRate(d time.Duration) Option {
	return func(s *StatusBar) {
		if d > 0 {
			s.interval = d
		}
	}
}

// WithStyle sets the visual style (StyleASCII or StyleEmoji).
func WithStyle(style Style) Option {
	return func(s *StatusBar) {
		s.style = style
	}
}

// WithColors sets the ANSI foreground and background color codes.
// Use the Fg* and Bg* constants defined in this package.
func WithColors(fg, bg int) Option {
	return func(s *StatusBar) {
		s.fgColor = fg
		s.bgColor = bg
	}
}

// WithWidth sets the fixed display width (in terminal columns) of the status
// line. Content shorter than width is padded with spaces; content wider than
// width is truncated. Calling this option disables automatic terminal-width
// detection. A value of 0 disables width enforcement entirely.
func WithWidth(w int) Option {
	return func(s *StatusBar) {
		if w >= 0 {
			s.width = w
			s.autoWidth = false
		}
	}
}

// New creates a StatusBar. Call Start to begin rendering.
// Default style: StyleASCII, cyan background (BgCyan), black foreground (FgBlack).
// By default, the width is automatically derived from the terminal on every render.
// Use WithWidth to override this behavior.
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
// Calling Start on an already-running StatusBar is a no-op.
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

// Stop halts the background goroutine and clears the status line.
// Calling Stop on a StatusBar that is not running is a no-op.
func (s *StatusBar) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopCh)

	// Render the final frame and move the cursor to a new line so the status
	// bar remains visible in the terminal and subsequent output is unaffected.
	s.mu.Lock()
	s.renderLocked()
	_, _ = fmt.Fprint(s.out, "\n")
	s.mu.Unlock()
}

// SetMessage sets the message shown on the status line.
// Pass an empty string to clear it.
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

// Writer returns an io.Writer that serializes writes with the status bar's
// internal mutex. Before each write it clears the current status line; after
// the write it redraws the status bar, so that log output does not corrupt
// the display.
//
// Pass this writer as the underlying output of the log console writer:
//
//	cw := console.New(console.WithWriter(bar.Writer()))
func (s *StatusBar) Writer() io.Writer {
	return &statusWriter{bar: s}
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

// renderLocked writes the status line. Must be called with s.mu held.
func (s *StatusBar) renderLocked() {
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

	// \033[K fills the rest of the line with the active background color,
	// then \033[0m resets all attributes.
	_, _ = fmt.Fprintf(s.out, "\r\033[%dm\033[%dm%s\033[K\033[0m",
		s.bgColor, s.fgColor, content)
}

// formatElapsed formats d as hh:mm:ss, omitting the hours group when zero.
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

// queryTerminalWidth probes the current terminal width using several methods
// to maximize compatibility across platforms and environments:
//
//  1. term.GetSize on stdout, stderr, and stdin (in that order) — covers most
//     native shells on Windows (cmd / PowerShell / Windows Terminal), macOS,
//     and Linux, as well as JetBrains IDE run windows that back their output
//     pane with a PTY.
//  2. The COLUMNS environment variable — a fallback that works in shells that
//     export it (bash, zsh, fish) and in JetBrains run configurations where
//     the PTY size is not surfaced via ioctl.
//
// Returns 0 when none of the probes succeed, signaling "no constraint".
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

// runeDisplayWidth returns the number of terminal columns a rune occupies.
// Emoji and other wide characters (East Asian wide/fullwidth) count as 2;
// control characters count as 0; everything else counts as 1.
func runeDisplayWidth(r rune) int {
	switch {
	case r < 0x20 || r == 0x7F:
		return 0
	case r >= 0x1100 && (r <= 0x115F || // Hangul Jamo
		r == 0x2329 || r == 0x232A ||
		(r >= 0x2E80 && r <= 0x303E) || // CJK Radicals … CJK Symbols
		(r >= 0x3040 && r <= 0x33FF) || // Hiragana … CJK Compatibility
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Ext-A
		(r >= 0x4E00 && r <= 0xA4CF) || // CJK Unified … Yi
		(r >= 0xA960 && r <= 0xA97F) || // Hangul Jamo Ext-A
		(r >= 0xAC00 && r <= 0xD7FF) || // Hangul Syllables
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0xFE10 && r <= 0xFE1F) || // Vertical Forms
		(r >= 0xFE30 && r <= 0xFE6F) || // CJK Compatibility Forms
		(r >= 0xFF01 && r <= 0xFF60) || // Fullwidth Latin
		(r >= 0xFFE0 && r <= 0xFFE6) || // Fullwidth Signs
		(r >= 0x1B000 && r <= 0x1B0FF) || // Kana Supplement
		(r >= 0x1F004 && r <= 0x1F0CF) || // Mahjong / Playing cards
		(r >= 0x1F300 && r <= 0x1F9FF) || // Misc Symbols / Emoji
		(r >= 0x20000 && r <= 0x2FFFD) || // CJK Ext-B … Supp Ideo
		(r >= 0x30000 && r <= 0x3FFFD)):
		return 2
	default:
		return 1
	}
}

// displayWidth returns the total number of terminal columns the string occupies.
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		w += runeDisplayWidth(r)
	}
	return w
}

// fitWidth truncates or right-pads s so that its display width equals exactly n.
func fitWidth(s string, n int) string {
	cur := 0
	for i, r := range s {
		rw := runeDisplayWidth(r)
		if cur+rw > n {
			// Truncate: pad with spaces to reach exactly n if a wide char
			// would overshoot.
			result := s[:i]
			for cur < n {
				result += " "
				cur++
			}
			return result
		}
		cur += rw
	}
	// Pad with spaces if content is shorter than n.
	for cur < n {
		s += " "
		cur++
	}
	return s
}

// statusWriter is an io.Writer that clears and redraws the status bar around
// every Write call, ensuring log lines and the status bar do not overwrite
// each other.
type statusWriter struct {
	bar *StatusBar
}

func (w *statusWriter) Write(p []byte) (n int, err error) {
	w.bar.mu.Lock()
	defer w.bar.mu.Unlock()

	// Clear the current status line before writing the log content.
	_, _ = fmt.Fprint(w.bar.out, "\r\033[K")

	n, err = w.bar.out.Write(p)
	if err != nil {
		return
	}

	// Redraw the status bar on the same line (log output already ended with \n).
	if w.bar.running {
		w.bar.renderLocked()
	}
	return
}
