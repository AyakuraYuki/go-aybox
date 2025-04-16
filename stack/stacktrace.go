/*
Inspired by samber/oops repo
-> https://github.com/samber/oops/blob/main/stacktrace.go
-> MIT License

Inspired by palantir/Stacktrace repo
-> https://github.com/palantir/stacktrace/blob/master/stacktrace.go
-> Apache 2.0 LICENSE
*/

package stack

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type fake struct{}

var (
	StackTraceMaxDepth = 10

	packageName = reflect.TypeOf(fake{}).PkgPath()
)

type stacktraceFrame struct {
	pc      uintptr
	file    string
	rawFile string
	fn      string
	line    int
}

func (frame *stacktraceFrame) String() string {
	currentFrame := fmt.Sprintf("%v:%v", frame.file, frame.line)
	if frame.fn != "" {
		currentFrame = fmt.Sprintf("%v:%v %v()", frame.file, frame.line, frame.fn)
	}
	return currentFrame
}

type Stacktrace struct {
	span   string
	frames []stacktraceFrame
}

func (st *Stacktrace) Error() string {
	return st.String("")
}

func (st *Stacktrace) String(deepestFrame string) string {
	var str string

	newline := func() {
		if str != "" && !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
	}

	for _, frame := range st.frames {
		if frame.file != "" {
			currentFrame := frame.String()
			if currentFrame == deepestFrame {
				break
			}
			newline()
			str += "  --- at " + currentFrame
		}
	}

	return str
}

func (st *Stacktrace) Source() (string, []string) {
	if len(st.frames) == 0 {
		return "", make([]string, 0)
	}
	firstFrame := st.frames[0]
	header := firstFrame.String()
	body := getSourceFromFrame(firstFrame)
	return header, body
}

func (st *Stacktrace) Sources() string {
	header, body := st.Source()
	if header == "" {
		header = "Thrown:"
	}
	str := ""
	str += header
	str += "\n"
	str += strings.Join(body, "\n")
	return str
}

func NewStacktrace(span ...string) *Stacktrace {
	sp := "stack"
	if len(span) > 0 && strings.TrimSpace(span[0]) != "" {
		sp = strings.TrimSpace(span[0])
	}

	frames := make([]stacktraceFrame, 0)

	// loop until we got StackTraceMaxDepth frames or run out of frames,
	// frames from this package are skipped
	for i := 1; len(frames) < StackTraceMaxDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		rawFile := file
		file = RemoveGoPath(file)

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		fn := shortFuncName(f)

		packageNameExamples := packageName + "/examples/"

		isGoPkg := len(runtime.GOROOT()) > 0 && strings.Contains(file, runtime.GOROOT()) // skip frames in GOROOT if it's set
		isThisPkg := strings.Contains(file, packageName)                                 // skip frames in this package
		isExamplePkg := strings.Contains(file, packageNameExamples)                      // do not skip frames in this package examples
		isTestPkg := strings.Contains(file, "_test.go")                                  // do not skip frames in tests

		if !isGoPkg && (!isThisPkg || isExamplePkg || isTestPkg) {
			frames = append(frames, stacktraceFrame{
				pc:      pc,
				file:    file,
				rawFile: rawFile,
				fn:      fn,
				line:    line,
			})
		}
	}

	return &Stacktrace{
		span:   sp,
		frames: frames,
	}
}

func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/example/proj/package.FuncName"
	// - "github.com/example/proj/package.Receiver.MethodName"
	// - "github.com/example/proj/package.(*PtrReceiver).MethodName"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}
