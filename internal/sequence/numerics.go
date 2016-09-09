package sequence

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/draw"
)

// SequentialNumerics provides sequential (column-wise) rendering calculations
type SequenceNumerics interface {
	Sequence() []base.PixelMember
}

func ImageSequence(sn SequenceNumerics, context draw.DrawingContext) {
	members := sn.Sequence()
	for _, point := range members {
		draw.DrawPoint(context, point)
	}
}

func Capture(sn SequenceNumerics) []base.PixelMember {
	return sn.Sequence()
}
