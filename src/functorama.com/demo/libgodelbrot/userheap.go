package libgodelbrot

type HeapItem interface{}


// A user managed heap
type UserHeap interface {
    Get(index int) HeapItem
    Add(item HeapItem)
}

type BaseHeap struct {
    index int
}

func (base *BaseHeap) Grow(heap UserHeap, item HeapItem) HeapItem {
    heap.Add(item)
    index := heap.index
    heap.index++
    return &heap.zone[index]
}