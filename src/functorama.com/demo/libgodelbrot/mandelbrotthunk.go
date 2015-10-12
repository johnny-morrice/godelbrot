package libgodelbrot

type MandelbrotThunk interface {
    MandelbrotMember
    // Note the thunk is done
    MarkComplete()
    // True if this thunk has beeen executed
    Done() bool
}

// Evaluated the thunk
func EvalThunk(thunk MandelbrotThunk) {
    if !thunk.Done() {
        thunk.Mandelbrot(Native.IterateLimit, Native.DivergeLimit)
        thunk.MarkComplete()
    }
}

type BaseThunk struct {
    evaled bool
}

func (thunk *BaseThunk) Done() bool {
    return thunk.evaled
}

func (thunk *BaseThunk) MarkComplete() {
    thunk.evaled
}