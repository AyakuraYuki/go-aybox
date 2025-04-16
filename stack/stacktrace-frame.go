package stack

import "fmt"

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
