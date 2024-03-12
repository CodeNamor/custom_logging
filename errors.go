package logging

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	maxStackLength = 50
	skipLevelThree = 3
)

type Error struct {
	Err        error
	StackTrace string
}

func (m Error) Error() string {
	return m.Err.Error() + m.StackTrace
}

func getStackTrace() string {
	stackBuf := make([]uintptr, maxStackLength)
	length := runtime.Callers(skipLevelThree, stackBuf)
	stack := stackBuf[:length]
	trace := ""

	frames := runtime.CallersFrames(stack)

	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace += fmt.Sprintf(" File: [%s], Line: [%d]. Function: [%s] | ", frame.File, frame.Line, frame.Function)
		}

		if !more {
			break
		}
	}

	return trace
}
