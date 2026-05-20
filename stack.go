package apperror

import (
	"runtime"
	"strconv"
)

// maxStackDepth caps how many frames newAppError records. Deep enough to
// reach the inbound boundary from any realistic domain/application call
// site, shallow enough that runaway recursion can't bloat a single error.
const maxStackDepth = 32

// StackTrace is the sequence of program counters captured at the point an
// AppError was constructed, innermost (origin) frame first.
//
// It stores raw PCs rather than resolved frames on purpose: symbolizing a
// stack (file/line/function lookup) is comparatively expensive, and most
// errors are created on paths that never log a trace. Resolution is
// deferred to Frames / Strings, which the logging boundary calls only for
// the errors it actually decides to record.
type StackTrace []uintptr

// callers records the call stack at the point an AppError is built. skip is
// the number of leading frames to drop so the first recorded frame is the
// caller that constructed the error, not this package's own construction
// frames. It is marked noinline so the frame layout skip relies on stays
// fixed regardless of the optimizer.
//
//go:noinline
func callers(skip int) StackTrace {
	var pcs [maxStackDepth]uintptr
	n := runtime.Callers(skip, pcs[:])
	st := make(StackTrace, n)
	copy(st, pcs[:n])
	return st
}

// Frames resolves the recorded program counters to runtime.Frame values,
// innermost first, expanding any inlined calls. Returns nil for an empty
// trace.
func (st StackTrace) Frames() []runtime.Frame {
	if len(st) == 0 {
		return nil
	}
	cf := runtime.CallersFrames(st)
	out := make([]runtime.Frame, 0, len(st))
	for {
		f, more := cf.Next()
		out = append(out, f)
		if !more {
			break
		}
	}
	return out
}

// Strings renders the trace as one "file:line function" entry per frame,
// innermost first. Each entry is single-line by construction, so the result
// is safe to attach as a structured log field without tripping the
// no-newline-in-messages rule (see error-handling guidance). Returns nil
// for an empty trace.
func (st StackTrace) Strings() []string {
	frames := st.Frames()
	if len(frames) == 0 {
		return nil
	}
	out := make([]string, len(frames))
	for i, f := range frames {
		out[i] = f.File + ":" + strconv.Itoa(f.Line) + " " + f.Function
	}
	return out
}
