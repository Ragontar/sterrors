package sterrors

import (
	"fmt"
	"runtime"
)

// STACK_TRACE_DEPTH определяет максимальную глубину возвращаемого стектрейса.
const STACK_TRACE_DEPTH = 30

type StackTraceError interface {
	Trace() string
	Error() string
}

type StackTracer interface {
	Trace() string
}

type stackTrace struct {
	trace string
}

// Trace вовзращает стектрейс в виде строки. Если не использовался WithStackTrace при создании ошибки, то вернется
// пустая строка
func (t stackTrace) Trace() string {
	return t.trace
}

func newStackTrace() stackTrace {
	pc := make([]uintptr, STACK_TRACE_DEPTH)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	var trace string
	for {
		frame, more := frames.Next()

		trace += fmt.Sprintf("%s:%d\n", frame.File, frame.Line)
		trace += fmt.Sprintf("\t%s\n", frame.Function)

		if !more {
			break
		}
	}

	return stackTrace{trace: trace}
}
