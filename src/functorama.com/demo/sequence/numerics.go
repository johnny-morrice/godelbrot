package sequence

import (
	"image"
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

// SequentialNumerics provides sequential (column-wise) rendering calculations
type SequenceNumerics interface {
    MandelbrotSequence(iterateLimit uint8)
    ImageDrawSequencer(context draw.DrawingContext)
    MemberCaptureSequencer()
    CapturedMembers() []base.PixelMember
    SubImage(rect image.Rectangle)
}

func Capture(sequence SequenceNumerics, iterateLimit uint8) []base.PixelMember {
    sequence.MemberCaptureSequencer()
    sequence.MandelbrotSequence(iterateLimit)
    return sequence.CapturedMembers()
}