package nativesequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/nativebase"
	"functorama.com/demo/sequence"
)

type NativeSequenceNumerics struct {
	nativebase.NativeBaseNumerics
	lastArea int
}

// Check we implement interface
var _ sequence.SequenceNumerics = (*NativeSequenceNumerics)(nil)

func Make(app nativebase.RenderApplication) NativeSequenceNumerics {
	return NativeSequenceNumerics{
		NativeBaseNumerics: nativebase.Make(app),
	}
}

func (nsn *NativeSequenceNumerics) Sequence() []base.PixelMember {
	imageLeft, imageTop := nsn.PictureMin()
	imageRight, imageBottom := nsn.PictureMax()
	rUnit, iUnit := nsn.PixelSize()
	sqrtDl := nsn.SqrtDivergeLimit
	iterlim := nsn.IterateLimit

	area := (imageRight - imageLeft) * (imageBottom - imageTop)
	out := make([]base.PixelMember, area)

	count := 0
	x := nsn.RealMin
	for i := imageLeft; i < imageRight; i++ {
		y := nsn.ImagMax
		for j := imageTop; j < imageBottom; j++ {
			member := nativebase.NativeMandelbrotMember{
				C: complex(x, y),
				SqrtDivergeLimit: sqrtDl,
			}
			member.Mandelbrot(iterlim)
			out[count] = base.PixelMember{I: i, J: j, Member: member.BaseMandelbrot}
			y -= iUnit
			count++
		}
		x += rUnit
	}
	return out
}