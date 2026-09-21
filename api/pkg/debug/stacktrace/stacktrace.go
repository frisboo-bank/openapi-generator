package stacktrace

import (
	"fmt"
	"runtime"
	"strings"
)

func Capture(skip int) string {
	const depth = 32
	var pcs [depth]uintptr

	n := runtime.Callers(2+skip, pcs[:])
	if n == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:n])

	var sb strings.Builder

	sb.WriteString("goroutine … [running]:\n")
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&sb, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	return sb.String()
}
