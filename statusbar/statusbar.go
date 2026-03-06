package statusbar

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// StatusBar is a non-blocking bottom status bar running in a coroutine manner.
// It dynamically displays the current time, the duration the program has been
// running, and the name of the method currently being executed.
type StatusBar struct {
	startTime time.Time
	taskName  string
	mu        sync.RWMutex
	stopCh    chan struct{}
	doneCh    chan struct{}
	running   bool
	refreshMs int
	style     Style
	// The number of lines occupied by the status bar last time, used for
	// clearing.
	savedLines int
}

// Style defines the visual style of the status bar.
type Style struct {
	TimeIcon    string
	ElapsedIcon string
	TaskIcon    string
	Separator   string
	BarChar     string
	LeftCap     string
	RightCap    string
}

// DefaultStyle gives the default visual style.
func DefaultStyle() Style {
	return Style{
		TimeIcon:    "🕐",
		ElapsedIcon: "⏱",
		TaskIcon:    "⚡",
		Separator:   " │ ",
		BarChar:     "─",
		LeftCap:     "┤",
		RightCap:    "├",
	}
}

// MinimalStyle gives the minimalist ASCII visual style.
func MinimalStyle() Style {
	return Style{
		TimeIcon:    "[T]",
		ElapsedIcon: "[E]",
		TaskIcon:    "[>]",
		Separator:   " | ",
		BarChar:     "-",
		LeftCap:     "|",
		RightCap:    "|",
	}
}

type Option func(*StatusBar)

// WithRefreshRate set refresh interval (in milliseconds).
func WithRefreshRate(ms int) Option {
	return func(sb *StatusBar) {
		if ms > 0 {
			sb.refreshMs = ms
		}
	}
}

// WithStyle set visual style.
func WithStyle(s Style) Option {
	return func(sb *StatusBar) {
		sb.style = s
	}
}

func New() *StatusBar {
	return &StatusBar{
		startTime: time.Now(),
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
		refreshMs: 200,
		style:     DefaultStyle(),
	}
}

// NewWithOptions creates a StatusBar with custom options.
func NewWithOptions(opts ...Option) *StatusBar {
	sb := New()
	for _, opt := range opts {
		opt(sb)
	}
	return sb
}

// SetTask registers the name of the currently executing method/task
// (thread-safe).
func (sb *StatusBar) SetTask(name string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.taskName = name
}

// ClearTask clears the current task name
func (sb *StatusBar) ClearTask() {
	sb.SetTask("")
}

// Start the status bar in an asynchronous coroutine manner.
func (sb *StatusBar) Start() {
	sb.mu.Lock()
	if sb.running {
		sb.mu.Unlock()
		return
	}
	sb.running = true
	sb.startTime = time.Now()
	sb.mu.Unlock()

	// hide cursor
	_, _ = fmt.Fprint(os.Stderr, "\033[?25l")

	go sb.loop()
}

// Stop the status bar and clear the terminal status.
func (sb *StatusBar) Stop() {
	sb.mu.Lock()
	if !sb.running {
		sb.mu.Unlock()
		return
	}
	sb.running = false
	sb.mu.Unlock()

	close(sb.stopCh)
	<-sb.doneCh

	// clean the status
	sb.clearBar()
	_, _ = fmt.Fprint(os.Stderr, "\033[?25h")
}

func (sb *StatusBar) loop() {
	defer close(sb.doneCh)

	ticker := time.NewTicker(time.Duration(sb.refreshMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sb.stopCh:
			return
		case <-ticker.C:
			sb.render()
		}
	}
}

func (sb *StatusBar) getTermWidth() int {
	w, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

func (sb *StatusBar) render() {
	sb.mu.RLock()
	task := sb.taskName
	start := sb.startTime
	sb.mu.RUnlock()

	now := time.Now()
	elapsed := now.Sub(start)

	// Construct the various sections of content.
	timePart := fmt.Sprintf("%s %s", sb.style.TimeIcon, now.Format("15:04:05"))
	elapsedPart := fmt.Sprintf("%s %s", sb.style.ElapsedIcon, formatDuration(elapsed))

	var taskPart string
	if task != "" {
		taskPart = fmt.Sprintf("%s %s", sb.style.TaskIcon, task)
	}

	// Assembly content
	content := timePart + sb.style.Separator + elapsedPart
	if taskPart != "" {
		content = content + sb.style.Separator + taskPart
	}

	width := sb.getTermWidth()

	// Constructing separator lines and status lines
	barLine := sb.style.RightCap + strings.Repeat(sb.style.BarChar, max(0, width-2)) + sb.style.LeftCap
	statusLine := padOrTruncate(content, width)

	// First clear the previous status bar.
	sb.clearBar()

	// Save cursor → Jump to bottom → Write to status bar → Restore cursor
	output := "\0337"                                // Save cursor position
	output += fmt.Sprintf("\033[%d;1H", 9999)        // Move to the bottom line of the terminal (use large numbers)
	output += "\033[1A"                              // Move up one line (leave two lines for the status bar)
	output += "\033[2K"                              // Clear the current line
	output += "\033[36m" + barLine + "\033[0m"       // Blue separator line
	output += "\n\033[2K"                            // New line and clear
	output += "\033[97;46m" + statusLine + "\033[0m" // Text in black on a blue background
	output += "\0338"                                // Restore cursor position

	_, _ = fmt.Fprint(os.Stderr, output)
	sb.savedLines = 2
}

func (sb *StatusBar) clearBar() {
	if sb.savedLines <= 0 {
		return
	}
	output := "\0337"
	output += fmt.Sprintf("\033[%d;1H", 9999)
	for i := 0; i < sb.savedLines; i++ {
		output += fmt.Sprintf("\033[%dA\033[2K", 1)
	}
	output += "\0338"
	_, _ = fmt.Fprint(os.Stderr, output)
	sb.savedLines = 0
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func padOrTruncate(s string, width int) string {
	// Calculate the visible character width (approximately estimated, emojis
	// take up 2 spaces)
	visLen := visibleLength(s)
	if visLen >= width {
		return truncateVisible(s, width)
	}
	return s + strings.Repeat(" ", width-visLen)
}

func visibleLength(s string) int {
	n := 0
	for _, r := range s {
		if r > 0x1F00 { // Roughly judge wide characters/emojis
			n += 2
		} else {
			n++
		}
	}
	return n
}

func truncateVisible(s string, maxWidth int) string {
	n := 0
	for i, r := range s {
		w := 1
		if r > 0x1F00 {
			w = 2
		}
		if n+w > maxWidth {
			return s[:i]
		}
		n += w
	}
	return s
}
