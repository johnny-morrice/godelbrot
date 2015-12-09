package nativesequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/nativebase"
	"functorama.com/demo/sequence"
)

type NativeSequenceNumerics struct {
	nativebase.NativeBaseNumerics
	area int
}

// Check we implement interface
var _ sequence.SequenceNumerics = (*NativeSequenceNumerics)(nil)

func Make(app nativebase.RenderApplication) NativeSequenceNumerics {
	w, h := app.PictureDimensions()
	return NativeSequenceNumerics{
		NativeBaseNumerics: nativebase.Make(app),
		area: int(w * h),
	}
}

func (nsn *NativeSequenceNumerics) Area() int {
	return nsn.area
}

func (nsn *NativeSequenceNumerics) Sequence(iterLimit uint8) <-chan base.PixelMember {
	imageLeft, imageTop := nsn.PictureMin()
	imageRight, imageBottom := nsn.PictureMax()
	rUnit, iUnit := nsn.PixelSize()
	sqrtDl := nsn.SqrtDivergeLimit

	out := make(chan base.PixelMember)

	// This goroutine will exit once all members have been read out
	go func() {
		x := nsn.RealMin
		for i := imageLeft; i < imageRight; i++ {
			y := nsn.ImagMax
			for j := imageTop; j < imageBottom; j++ {
				member := nativebase.NativeMandelbrotMember{
					C: complex(x, y),
					SqrtDivergeLimit: sqrtDl,
				}
				member.Mandelbrot(iterLimit)
				out<- base.PixelMember{I: i, J: j, Member: member.BaseMandelbrot}
				y -= iUnit
			}
			x += rUnit
		}
		close(out)
	}()

	return out
}