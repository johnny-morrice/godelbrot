package libgodelbrot

type EscapePoint struct {
	evaluated  bool
	c          complex128
	membership MandelbrotMember
}

func NewEscapePointReals(r float64, i float64) *EscapePoint {
	return NewEscapePoint(complex(r, i))
}

func NewEscapePoint(c complex128) *EscapePoint {
	return &EscapePoint{
		evaluated: false,
		c:         c,
	}
}

type EscapePointHeap struct {
	zone  []EscapePoint
	index int
}

func NewEscapePointHeap(size uint) *EscapePointHeap {
	return &EscapePointHeap{
		zone:  make([]EscapePoint, 0, int(size)),
		index: 0,
	}
}

func (heap *EscapePointHeap) EscapePoint(x float64, y float64) *EscapePoint {
	heap.zone = append(heap.zone, EscapePoint{evaluated: false, c: complex(x, y)})
	index := heap.index
	heap.index++
	return &heap.zone[index]
}
