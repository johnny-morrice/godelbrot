package libgodelbrot

type RenderConfigHeap struct {
	zone  []RenderConfig
	base  RenderConfig
	index int
}

func NewRenderConfigHeap(baseConfig *RenderConfig, size uint) *RenderConfigHeap {
	heap := &RenderConfigHeap{
		zone:  make([]RenderConfig, 0, size),
		base:  *baseConfig,
		index: 0,
	}
	heap.base.Frame = CornerFrame
	return heap
}

func (heap *RenderConfigHeap) Subconfig(region Region) *RenderConfig {
	heap.zone = append(heap.zone, heap.base)
	config := &heap.zone[heap.index]
	heap.index++

	left, top := heap.base.PlaneToPixel(region.topLeft.c)
	right, bottom := heap.base.PlaneToPixel(region.bottomRight.c)
	config.Width = uint(right - left)
	config.Height = uint(bottom - top)
	config.ImageLeft = left
	config.ImageTop = top
	config.TopLeft = region.topLeft.c
	config.BottomRight = region.bottomRight.c
	return config
}
