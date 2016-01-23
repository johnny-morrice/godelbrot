package sequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

// SequentialNumerics provides sequential (column-wise) rendering calculations
type SequenceNumerics interface {
    Sequence() <-chan base.PixelMember
    Area() int
}

func ImageSequence(sn SequenceNumerics, context draw.DrawingContext) {
    members := sn.Sequence()
    for point := range members {
        draw.DrawPoint(context, point)
    }
}

func Capture(sn SequenceNumerics) []base.PixelMember {
    members := sn.Sequence()
    out := make([]base.PixelMember, sn.Area())
    i := 0
    for point := range members {
        out[i] = point
        i++
    }
    return out
}