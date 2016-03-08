package nativesequence

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/nativebase"
	"github.com/johnny-morrice/godelbrot/internal/sequence"
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
			member := nativebase.NativeEscapeValue{
				C: complex(x, y),
				SqrtDivergeLimit: sqrtDl,
			}
			member.Mandelbrot(iterlim)
			out[count] = base.PixelMember{I: i, J: j, Member: member.EscapeValue}
			y -= iUnit
			count++
		}
		x += rUnit
	}
	return out
}