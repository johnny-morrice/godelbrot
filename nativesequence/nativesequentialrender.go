package nativesequence

import (
	"github.com/johnny-morrice/godelbrot/base"
	"github.com/johnny-morrice/godelbrot/nativebase"
	"github.com/johnny-morrice/godelbrot/sequence"
)

type NativeSequenceNumerics struct {
	nativebase.NativeBaseNumerics
}

// Check we implement interface
var _ sequence.SequenceNumerics = (*NativeSequenceNumerics)(nil)

func Make(app nativebase.RenderApplication) NativeSequenceNumerics {
	return NativeSequenceNumerics{
		NativeBaseNumerics: nativebase.Make(app),
	}
}

func (nsn *NativeSequenceNumerics) Sequence() []base.PixelMember {
	ileft, itop := nsn.PictureMin()
	iright, ibott := nsn.PictureMax()
	rUnit, iUnit := nsn.PixelSize()
	sqrtDl := nsn.SqrtDivergeLimit
	iterlim := nsn.IterateLimit

	area := (iright - ileft) * (ibott - itop)
	out := make([]base.PixelMember, area)

	count := 0
	x := nsn.RealMin
	for i := ileft; i < iright; i++ {
		y := nsn.ImagMax
		for j := itop; j < ibott; j++ {
			member := nativebase.NativeMandelbrotMember{
				C: complex(x, y),
				SqrtDivergeLimit: sqrtDl,
			}
			member.Mandelbrot(iterlim)
			out[count] = base.PixelMember{I: i, J: j, Member: member.MandelbrotMember}
			y -= iUnit
			count++
		}
		x += rUnit
	}
	return out
}