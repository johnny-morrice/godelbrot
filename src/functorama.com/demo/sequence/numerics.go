package sequence

import (
	"image"
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

type SequenceNumericsFactory interface {
    Build() SequenceNumerics
}

// SequentialNumerics provides sequential (column-wise) rendering calculations
type SequenceNumerics interface {
    MandelbrotSequence(iterateLimit uint8)
    ImageDrawSequencer(context draw.DrawingContext)
    MemberCaptureSequencer()
    CapturedMembers() []base.PixelMember
    SubImage(rect image.Rectangle)
}