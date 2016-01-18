package sequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

// SequentialNumerics provides sequential (column-wise) rendering calculations
type SequenceNumerics interface {
    Sequence(iterateLimit uint8) <-chan base.PixelMember
    Area() int
}

func ImageSequence(sn SequenceNumerics, iterateLimit uint8, context draw.DrawingContext) {
    members := sn.Sequence(iterateLimit)
    for point := range members {
        draw.DrawPoint(context, point)
    }
}

func Capture(sn SequenceNumerics, iterateLimit uint8) []base.PixelMember {
    members := sn.Sequence(iterateLimit)
    out := make([]base.PixelMember, sn.Area())
    i := 0
    for point := range members {
        out[i] = point
        i++
    }
    return out
}