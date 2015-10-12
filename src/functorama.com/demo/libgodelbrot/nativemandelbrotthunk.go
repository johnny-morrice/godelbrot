package libgodelbrot

// NativeMandelbrotMember that tracks if it has been evaluated yet
type NativeMandelbrotThunk struct {
	BaseThunk
	NativeMandelbrotMember
}

func NewNativeMandelbrotThunkReals(r float64, i float64) *NativeMandelbrotThunk {
	return NewNativeMandelbrotThunk(complex(r, i))
}

func NewNativeMandelbrotThunk(c complex128) *NativeMandelbrotThunk {
	return &NativeMandelbrotThunk{
		NativeMandelbrotMember: NativeMandelbrotMember{
			C: c
		}
	}
}

// Heap for creating many NativeMandelbrotThunks
type NativeMandelbrotThunkHeap struct {
	zone  []NativeMandelbrotThunk
	BaseHeap
}

func NewNativeMandelbrotThunkHeap(size uint) *NativeMandelbrotThunkHeap {
	return &NativeMandelbrotThunkHeap{
		zone:  make([]NativeMandelbrotThunk, 0, int(size)),
	}
}

func (heap *NativeMandelbrotThunkHeap) Get(index int) HeapItem {
	return heap.zone[index]
}

func (heap *NativeMandelbrotThunkHeap) Add(thunk HeapItem) {
	th := thunk.(NativeMandelbrotThunk) 
	heap.zone = append(heap.zone, th)
}

func (heap *NativeMandelbrotThunkHeap) NativeMandelbrotThunk(x float64, y float64) *NativeMandelbrotThunk {
	thunk := NativeMandelbrotThunk{
		NativeMandelbrotMember: NativeMandelbrotMember{
			C: c
		}
	}
	item := heap.BaseHeap.Grow(heap, thunk)
	return &item.(NativeMandelbrotThunk)
}
