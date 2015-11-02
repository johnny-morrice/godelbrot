package nativesequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
	"functorama.com/demo/nativebase"
)

type NativeSequenceNumerics struct {
	nativebase.NativeBaseNumerics
	sequencer func(i int, j int, member nativebase.NativeMandelbrotMember)
	members   []base.PixelMember
}

func CreateNativeSequenceNumerics(base nativebase.NativeBaseNumerics) NativeSequenceNumerics {
	return NativeSequenceNumerics{NativeBaseNumerics: base}
}

func (native *NativeSequenceNumerics) MandelbrotSequence(iterLimit uint8) {
	topLeft := native.PlaneTopLeft()

	imageLeft, imageTop := native.PictureMin()
	imageRight, imageBottom := native.PictureMax()
	rUnit, iUnit := native.PixelSize()
	sqrtDl := native.SqrtDivergeLimit

	x := real(topLeft)
	for i := imageLeft; i < imageRight; i++ {
		y := imag(topLeft)
		for j := imageTop; j < imageBottom; j++ {
			member := nativebase.NativeMandelbrotMember{
				C: complex(x, y),
				SqrtDivergeLimit: sqrtDl,
			}
			member.Mandelbrot(iterLimit)
			native.sequencer(i, j, member)
			y -= iUnit
		}
		x += rUnit
	}
}

func (native *NativeSequenceNumerics) ImageDrawSequencer(context draw.DrawingContext) {
	native.sequencer = func(i, j int, member nativebase.NativeMandelbrotMember) {
		draw.DrawPoint(context, base.PixelMember{i, j, &member})
	}
}

func (native *NativeSequenceNumerics) MemberCaptureSequencer() {
	native.sequencer = func(i, j int, member nativebase.NativeMandelbrotMember) {
		native.members = append(native.members, base.PixelMember{
			I:      i,
			J:      j,
			Member: &member,
		})
	}
}

func (native *NativeSequenceNumerics) CapturedMembers() []base.PixelMember {
	return native.members
}
