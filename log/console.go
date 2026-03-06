package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog"
)

var _ zerolog.LevelWriter = (*ConsoleWriter)(nil)

const consoleFormatJSON = "json"

const (
	cReset    = 0
	cBold     = 1
	cRed      = 31
	cGreen    = 32
	cYellow   = 33
	cBlue     = 34
	cMagenta  = 35
	cCyan     = 36
	cGray     = 37
	cDarkGray = 90
)

var consoleBufPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 100))
	},
}

type ConsoleWriter struct {
	level   zerolog.Level
	format  string // empty or JSON
	noColor bool
	out     io.Writer
}

type ConsoleWriterOption func(*ConsoleWriter)

func NewConsoleWriter(opts ...ConsoleWriterOption) io.Writer {
	writer := &ConsoleWriter{
		level:   zerolog.DebugLevel,
		noColor: runInK8S(),
		out:     os.Stdout,
	}

	for _, opt := range opts {
		opt(writer)
	}

	if async {
		asyncWriter := NewAsyncWriter(writer.level, writer, diodes.NewManyToOne(1024, diodes.AlertFunc(func(missed int) {
			fmt.Printf("console writer dropped %d messages\n", missed)
		})), 1*time.Second)
		registerCloseFn(asyncWriter.Close)
		return asyncWriter
	}

	return writer
}

// Write data to writer.
func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	if w.format == consoleFormatJSON {
		_, _ = w.out.Write(p)
		n = len(p)
		return
	}

	var event map[string]any
	err = sonic.Unmarshal(p, &event)
	if err != nil {
		return
	}

	buf := consoleBufPool.Get().(*bytes.Buffer)
	defer consoleBufPool.Put(buf)

	var (
		level    = "?????"
		lvlColor = cReset
	)
	if l, ok := event[zerolog.LevelFieldName].(string); ok {
		level = l
		lvlColor = levelColor(l, w.noColor)
	}
	if _, ok := event[zerolog.TimestampFieldName]; ok {
		event[zerolog.TimestampFieldName] = time.Now().Format("2006-01-02 15:04:05.999999")
	}
	_, _ = fmt.Fprintf(buf, "%s | %-5s | %s |",
		colorize(event[zerolog.TimestampFieldName], cDarkGray, w.noColor),
		colorize(level, lvlColor, w.noColor),
		colorize(event[zerolog.MessageFieldName], cReset, w.noColor))

	fields := make([]string, 0)
	for field := range event {
		switch field {
		case zerolog.LevelFieldName, zerolog.TimestampFieldName, zerolog.MessageFieldName:
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		_, _ = fmt.Fprintf(buf, " %s=", colorize(field, cCyan, w.noColor))
		switch value := event[field].(type) {
		case string:
			if needsQuote(value) {
				buf.WriteString(strconv.Quote(value))
			} else {
				buf.WriteString(value)
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			_, _ = fmt.Fprint(buf, value)
		default:
			bs, err := sonic.Marshal(value)
			if err != nil {
				_, _ = fmt.Fprintf(buf, "[error: %v]", err)
			} else {
				_, _ = fmt.Fprint(buf, string(bs))
			}
		}
	}

	buf.WriteByte('\n')
	_, _ = buf.WriteTo(w.out)
	n = len(p)
	return
}

// WriteLevel writes data to writer with level info provided
func (w *ConsoleWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < w.level {
		return len(p), nil
	}
	return w.Write(p)
}

func colorize(s any, color int, noColor bool) string {
	if noColor {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, s)
}

func levelColor(level string, noColor bool) int {
	if noColor {
		return cReset
	}
	switch strings.ToLower(level) {
	case "debug":
		return cMagenta
	case "info":
		return cGreen
	case "warn":
		return cYellow
	case "error", "fatal", "panic":
		return cRed
	default:
		return cReset
	}
}

func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}

func WithConsoleLogLevel(level zerolog.Level) ConsoleWriterOption {
	return func(w *ConsoleWriter) {
		w.level = level
	}
}

func WithConsoleFormatJSON() ConsoleWriterOption {
	return func(w *ConsoleWriter) {
		w.format = consoleFormatJSON
	}
}

func WithConsoleNoColor() ConsoleWriterOption {
	return func(w *ConsoleWriter) {
		w.noColor = true
	}
}

func WithConsoleWriter(out io.Writer) ConsoleWriterOption {
	return func(w *ConsoleWriter) {
		if out != nil {
			w.out = out
		}
	}
}
