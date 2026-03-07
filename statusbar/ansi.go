package statusbar

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

// ANSI escape sequences for cursor control and screen management.
const (
	CursorUp         = "\033[A"
	CursorDown       = "\033[B"
	CursorRight      = "\033[C"
	CursorLeft       = "\033[D"
	CursorNextLine   = "\033[E"
	CursorPrevLine   = "\033[F"
	CursorToColumn   = "\033[%dG"    // format string; use with fmt.Sprintf
	CursorToPosition = "\033[%d;%dH" // format string: row, col (1-indexed)
	SaveCursorPos    = "\033[s"
	RestoreCursorPos = "\033[u"
	ClearLine        = "\033[K"
	ClearScreen      = "\033[2J"
	EnableScrolling  = "\033[r"   // reset DECSTBM scroll region to full terminal
	AutoWrapOff      = "\033[?7l" // disable terminal auto-wrap (DECRST 7)
	AutoWrapOn       = "\033[?7h" // enable terminal auto-wrap (DECSET 7)
)
