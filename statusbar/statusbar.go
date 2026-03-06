package statusbar

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"
)

const reservedLines = 2 // 状态栏占用的底部行数（分隔线 + 状态行）

// StatusBar 是一个非阻塞的底部状态栏，以协程方式运行。
// 使用 ANSI Scroll Region 技术将终端分为上方滚动区和底部固定区，
// 确保主程序的日志输出不会冲掉状态栏。
type StatusBar struct {
	startTime  time.Time
	taskName   string
	mu         sync.RWMutex
	writeMu    sync.Mutex // serializes all writes to out (render + external writers)
	stopCh     chan struct{}
	doneCh     chan struct{}
	running    bool
	refreshMs  int
	style      Style
	termHeight int
	termWidth  int
	out        io.Writer
}

// Style 定义状态栏的视觉样式
type Style struct {
	TimeIcon    string
	ElapsedIcon string
	TaskIcon    string
	Separator   string
	BarChar     string
	LeftCap     string
	RightCap    string
}

// DefaultStyle 返回默认样式
func DefaultStyle() Style {
	return Style{
		TimeIcon:    " 🕐",
		ElapsedIcon: "⏱",
		TaskIcon:    "⚡",
		Separator:   " │ ",
		BarChar:     "─",
		LeftCap:     "┤",
		RightCap:    "├",
	}
}

// MinimalStyle 返回简约ASCII样式
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

// Option 是 StatusBar 的配置函数类型
type Option func(*StatusBar)

// WithRefreshRate 设置刷新间隔（毫秒）
func WithRefreshRate(ms int) Option {
	return func(sb *StatusBar) {
		if ms > 0 {
			sb.refreshMs = ms
		}
	}
}

// WithStyle 设置视觉样式
func WithStyle(s Style) Option {
	return func(sb *StatusBar) {
		sb.style = s
	}
}

func WithWriter(w io.Writer) Option {
	return func(sb *StatusBar) {
		if w != nil {
			sb.out = w
		}
	}
}

// New 创建一个新的 StatusBar 实例
func New(opts ...Option) *StatusBar {
	bar := &StatusBar{
		startTime: time.Now(),
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
		refreshMs: 200,
		style:     DefaultStyle(),
		out:       os.Stdout,
	}
	for _, opt := range opts {
		opt(bar)
	}
	return bar
}

// Writer returns an io.Writer that serializes external writes with the status bar's
// render loop. Pass this as the output writer for any logger so that log lines and
// status bar drawing never interleave on the terminal.
//
// Example:
//
//	bar := statusbar.New()
//	bar.Start()
//	defer bar.Stop()
//
//	logger := log.New(log.WithOutput(log.NewConsoleWriter(log.WithConsoleWriter(bar.Writer()))))
func (sb *StatusBar) Writer() io.Writer {
	return &statusBarWriter{sb: sb}
}

// statusBarWriter is an io.Writer that acquires writeMu before writing,
// ensuring mutual exclusion with render().
type statusBarWriter struct {
	sb *StatusBar
}

func (w *statusBarWriter) Write(p []byte) (n int, err error) {
	w.sb.writeMu.Lock()
	defer w.sb.writeMu.Unlock()
	return w.sb.out.Write(p)
}

// write is the single internal path for all terminal output.
// It holds writeMu so that render() and external writers are mutually exclusive.
func (sb *StatusBar) write(s string) {
	sb.writeMu.Lock()
	_, _ = fmt.Fprint(sb.out, s)
	sb.writeMu.Unlock()
}

// SetTask 注册当前执行的方法/任务名称（线程安全）
func (sb *StatusBar) SetTask(name string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.taskName = name
}

// ClearTask 清除当前任务名称
func (sb *StatusBar) ClearTask() {
	sb.SetTask("")
}

// Start 以协程方式启动状态栏（非阻塞）
func (sb *StatusBar) Start() {
	sb.mu.Lock()
	if sb.running {
		sb.mu.Unlock()
		return
	}
	sb.running = true
	sb.startTime = time.Now()
	sb.mu.Unlock()

	// 获取终端尺寸并设置滚动区域
	sb.updateTermSize()
	sb.setupScrollRegion()

	// 隐藏光标
	sb.write("\033[?25l")

	go sb.loop()
	go sb.watchResize()
}

// Stop 停止状态栏并清理终端状态
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

	// 清除状态栏区域
	sb.clearStatusArea()
	// 恢复全屏滚动区域
	sb.resetScrollRegion()
	// 显示光标
	sb.write("\033[?25h")
}

// updateTermSize 获取并缓存终端尺寸
func (sb *StatusBar) updateTermSize() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 || h <= 0 {
		w, h = 80, 24
	}
	sb.mu.Lock()
	sb.termWidth = w
	sb.termHeight = h
	sb.mu.Unlock()
}

// setupScrollRegion 设置终端滚动区域，预留底部给状态栏。
//
// 核心原理：ANSI 转义序列 \033[1;Nr 将终端的可滚动区域
// 限定在第 1 行到第 N 行。第 N+1 行及以下成为"固定区域"，
// 主程序的所有输出（fmt.Println 等）只会在滚动区域内滚动，
// 不会影响底部的状态栏。
func (sb *StatusBar) setupScrollRegion() {
	sb.mu.RLock()
	h := sb.termHeight
	sb.mu.RUnlock()

	scrollEnd := h - reservedLines
	if scrollEnd < 1 {
		scrollEnd = 1
	}

	out := ""
	// 先输出空行，为底部状态栏腾出空间
	for i := 0; i < reservedLines; i++ {
		out += "\n"
	}
	// 设置滚动区域：第 1 行 ~ 第 scrollEnd 行
	out += fmt.Sprintf("\033[1;%dr", scrollEnd)
	// 将光标移回滚动区域内（底部）
	out += fmt.Sprintf("\033[%d;1H", scrollEnd)

	sb.write(out)
}

// resetScrollRegion 恢复终端为全屏滚动
func (sb *StatusBar) resetScrollRegion() {
	sb.mu.RLock()
	h := sb.termHeight
	sb.mu.RUnlock()

	sb.write(fmt.Sprintf("\033[1;%dr\033[%d;1H", h, h-reservedLines))
}

// clearStatusArea 清除底部状态栏区域
func (sb *StatusBar) clearStatusArea() {
	sb.mu.RLock()
	h := sb.termHeight
	sb.mu.RUnlock()

	out := "\0337" // 保存光标
	for i := 0; i < reservedLines; i++ {
		row := h - reservedLines + 1 + i
		out += fmt.Sprintf("\033[%d;1H\033[2K", row)
	}
	out += "\0338" // 恢复光标
	sb.write(out)
}

// watchResize 监听终端窗口大小变化 (SIGWINCH)，自动调整滚动区域
func (sb *StatusBar) watchResize() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	defer signal.Stop(sigCh)

	for {
		select {
		case <-sb.stopCh:
			return
		case <-sigCh:
			sb.updateTermSize()
			sb.setupScrollRegion()
			sb.render()
		}
	}
}

func (sb *StatusBar) loop() {
	defer close(sb.doneCh)

	// 立即渲染一次
	sb.render()

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

func (sb *StatusBar) render() {
	sb.mu.RLock()
	task := sb.taskName
	start := sb.startTime
	h := sb.termHeight
	w := sb.termWidth
	sb.mu.RUnlock()

	if h <= reservedLines || w <= 0 {
		return
	}

	now := time.Now()
	elapsed := now.Sub(start)

	// 构建各段内容
	timePart := fmt.Sprintf("%s %s", sb.style.TimeIcon, now.Format("15:04:05"))
	elapsedPart := fmt.Sprintf("%s %s", sb.style.ElapsedIcon, formatDuration(elapsed))

	var taskPart string
	if task != "" {
		taskPart = fmt.Sprintf("%s %s", sb.style.TaskIcon, task)
	}

	// 组装内容
	content := " " + timePart + sb.style.Separator + elapsedPart
	if taskPart != "" {
		content = content + sb.style.Separator + taskPart
	}

	// 分隔线行
	barLine := sb.style.RightCap + strings.Repeat(sb.style.BarChar, max(0, w-2)) + sb.style.LeftCap
	// 状态行（填充至整行宽度）
	statusLine := padOrTruncate(content, w)

	// 定位到固定区域绘制（滚动区域之外，不受滚动影响）
	barRow := h - 1 // 分隔线所在行
	statusRow := h  // 状态信息所在行

	out := "\0337" // 保存光标位置
	// 绘制分隔线
	out += fmt.Sprintf("\033[%d;1H\033[2K", barRow)
	out += "\033[36m" + barLine + "\033[0m"
	// 绘制状态行
	out += fmt.Sprintf("\033[%d;1H\033[2K", statusRow)
	out += "\033[30;46m" + statusLine + "\033[0m"
	out += "\0338" // 恢复光标到滚动区域内的原位置

	sb.write(out)
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
	visLen := visibleLength(s)
	if visLen >= width {
		return truncateVisible(s, width)
	}
	return s + strings.Repeat(" ", width-visLen+3)
}

func visibleLength(s string) int {
	n := 0
	for _, r := range s {
		if r > 0x1F00 {
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
