package statusbar

import (
	"io"
	"time"
)

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

// WithWidth sets the fixed display width in terminal columns. Content shorter
// than width is right-padded; wider content is truncated. Passing 0 disables
// width enforcement. Calling this option disables auto-width detection.
// Notice: Conflict to WithFixedToBottom.
func WithWidth(w int) Option {
	return func(s *StatusBar) {
		if w >= 0 {
			s.width = w
			s.autoWidth = false
		}
	}
}

// WithFixedToBottom pins the status bar to the last terminal line. On each
// render the cursor is saved, the bar is drawn at the bottom, and the cursor
// is restored, so normal output above it is unaffected. Auto-width is enabled
// automatically so the bar always fills the full terminal width.
// Notice: Conflict to WithWidth.
func WithFixedToBottom() Option {
	return func(s *StatusBar) {
		s.fixedToBottom = true
		s.autoWidth = true
		s.width = 0
	}
}
