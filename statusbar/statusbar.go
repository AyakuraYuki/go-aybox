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
	StyleASCII Style = iota // StyleASCII uses plain ASCII separators.
	StyleEmoji              // StyleEmoji uses emoji icons as field indicators.
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

// ANSI escape sequences for cursor positioning and scrolling.
const (
	CursorUp            = "\033[A"
	CursorDown          = "\033[B"
	CursorRight         = "\033[C"
	CursorLeft          = "\033[D"
	CursorNextLine      = "\033[E"
	CursorPrevLine      = "\033[F"
	CursorToColumn      = "\033[%dG"
	CursorToPosition    = "\033[%d;%dH"
	SaveCursorPos       = "\033[s"
	RestoreCursorPos    = "\033[u"
	ClearLine           = "\033[K"
	ClearScreen         = "\033[2J"
	EnableScrolling     = "\033[r"
	DisableScrolling    = "\033[?7l"
	EnableScrollingBack = "\033[?7h"
)

// StatusBar renders a single-line status display in the terminal showing the
// current time, program elapsed time, and an optional message. It runs in a
// background goroutine and serializes all terminal writes through a mutex so
// that log output interleaved via Writer() does not corrupt the display.
//
// Typical usage with the log module:
//
//	bar := statusbar.New(statusbar.WithFixedToBottom())
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
	width         int  // 0 means no constraint; overridden by autoWidth at render time
	autoWidth     bool // when true, query terminal width on each render
	fixedToBottom bool // status bar fixed to shell bottom
	height        int  // Terminal height, queried dynamically
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

// WithFixedToBottom enables the feature to keep the status bar fixed at the
// bottom of the terminal window.
// When fixed to bottom, auto-width should be enabled for proper fitting.
func WithFixedToBottom() Option {
	return func(s *StatusBar) {
		s.fixedToBottom = true
		s.autoWidth = true
		s.width = 0
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
		// Initialize new fields
		fixedToBottom: false,
		height:        0,
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

	// If fixed to bottom, we might want to reserve space or manage screen differently.
	// For simplicity here, we just start the loop which will handle positioning.
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

	// If fixed to bottom, restore normal scrolling and move cursor past our line
	if s.fixedToBottom {
		fmt.Fprintf(s.out, EnableScrollingBack) // Re-enable scrolling
		fmt.Fprintf(s.out, CursorNextLine)      // Move to next line below status bar
	}

	s.renderLocked()
	fmt.Fprint(s.out, "\n")
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

	// Prepare the final output string for the status bar
	statusBarOutput := fmt.Sprintf("\033[%dm\033[%dm%s\033[K\033[0m", s.bgColor, s.fgColor, content)

	if s.fixedToBottom {
		// --- Fixed-to-bottom logic ---

		// 1. Query terminal height
		h := queryTerminalHeight()
		if h <= 0 {
			// If we can't get height, fall back to standard rendering
			fmt.Fprintf(s.out, "\r%s", statusBarOutput)
			return
		}
		s.height = h // Store for potential future use if needed

		// 2. Disable general scrolling for the full screen
		// This prevents other output from pushing our status bar up.
		fmt.Fprint(s.out, DisableScrolling)

		// 3. Move cursor to the last line (height-th row)
		// ANSI escape sequence for cursor position is \033[row;colH
		// Row and col are 1-indexed.
		moveCmd := fmt.Sprintf(CursorToPosition, h, 1)
		fmt.Fprint(s.out, moveCmd)

		// 4. Print the status bar content on that fixed line
		fmt.Fprint(s.out, statusBarOutput)

	} else {
		// --- Standard logic (existing behavior) ---
		fmt.Fprintf(s.out, "\r%s", statusBarOutput)
	}
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

// queryTerminalHeight probes the current terminal height.
// It follows the same strategy as queryTerminalWidth but for rows (lines).
// 1. term.GetSize on stdout, stderr, and stdin.
// 2. The LINES environment variable.
// Returns 0 when none of the probes succeed.
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
// It delegates to golang.org/x/text/width which implements the Unicode East
// Asian Width standard: Wide and Fullwidth characters (including emoji and CJK)
// count as 2; control characters count as 0; everything else counts as 1.
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

	// --- Modified logic for fixed-to-bottom ---

	if w.bar.fixedToBottom && w.bar.running {
		// If fixed to bottom, we need to temporarily move away from the status bar line
		// to print the log message, then come back.

		// 1. Get the height to know where the status bar is
		h := queryTerminalHeight()
		if h <= 0 {
			// Fallback: behave like original code if height unknown
			_, _ = fmt.Fprint(w.bar.out, "\r"+ClearLine)
			n, err = w.bar.out.Write(p)
			if err == nil && w.bar.running {
				w.bar.renderLocked()
			}
			return
		}

		// 2. Save cursor position (optional, but good practice)
		_, _ = fmt.Fprint(w.bar.out, SaveCursorPos)

		// 3. Move cursor to a safe line above the status bar (e.g., h-1)
		// We assume other program output goes above our fixed line.
		// This step might be complex if we don't control other output.
		// A simpler approach: just print the log, knowing it will scroll normally
		// except for our fixed bottom line due to \033[r disabling scrolling for the whole screen.
		// The key is that renderLocked() always repositions and redraws the status bar.

		// Let's print the log normally, relying on the renderLocked logic to fix the bottom.
		// First, clear the current line where the log might appear (though this is less critical now)
		// Then print the log
		_, _ = fmt.Fprint(w.bar.out, "\r"+ClearLine) // Clear current line before log
		n, err = w.bar.out.Write(p)
		if err != nil {
			return
		}

		// 4. Restore the status bar by calling renderLocked again.
		// This will reposition the cursor to the bottom line and redraw it.
		w.bar.renderLocked()

		// 5. Restore the saved cursor position (optional, depends on desired behavior)
		// fmt.Fprint(w.bar.out, RestoreCursorPos) // This might be confusing if logs were printed.
		// For log writers, it's usually better to let the cursor stay after the log line.
		// So, we comment out the restore here. The status bar will be there anyway.
	} else {
		// Original logic for non-fixed status bar
		_, _ = fmt.Fprint(w.bar.out, "\r"+ClearLine) // Clear the current status line before writing log
		n, err = w.bar.out.Write(p)
		if err != nil {
			return
		}
		// Redraw the status bar on the same line (log output already ended with \n).
		if w.bar.running {
			w.bar.renderLocked()
		}
	}

	return
}
