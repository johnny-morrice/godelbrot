package nativesequence

import (
	"math"
	"functorama.com/demo/base"
	"functorama.com/demo/nativebase"
	"functorama.com/demo/draw"
)

type NativeSequentialNumerics struct {
	nativebase.NativeBaseNumerics
	sequencer func(i int, j int, member NativeMandelbrotMember)
	members   []base.PixelMember
}

func (native *NativeSequentialNumerics) MandelbrotSequence(iterLimit uint8) {
	topLeft := native.PlaneTopLeft()

	imageLeft, imageTop := native.PictureMin()
	imageRight, imageBottom := native.PictureMax()
	rUnit, iUnit := native.PixelSize()
	sqrtDl := math.Sqrt(native.DivergeLimit)

	x := real(topLeft)
	for i := imageLeft; i < imageRight; i++ {
		y := imag(topLeft)
		for j := imageTop; j < imageBottom; j++ {
			member := NativeMandelbrotMember{
				C: complex(x, y),
				SrqtDivergeLimit: sqrtDl,
			}
			member.Mandelbrot(iterLimit)
			native.sequencer(i, j, member)
			y -= iUnit
		}
		x += rUnit
	}
}

func (native *NativeSequentialNumerics) ImageDrawSequencer(context draw.DrawingContext) {
	native.sequencer = func(i, j int, member NativeMandelbrotMember) {
		draw.DrawPoint(context, base.PixelMember{i, j, &member})
	}
}

func (native *NativeSequentialNumerics) MemberCaptureSequencer() {
	native.sequencer = func(i, j int, member NativeMandelbrotMember) {
		native.members = append(native.members, base.PixelMember{
			I:      i,
			J:      j,
			Member: &member,
		})
	}
}

func (native *NativeSequentialNumerics) CapturedMembers() []base.PixelMember {
	return native.members
}
