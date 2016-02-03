package bigsequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/bigbase"
	"functorama.com/demo/sequence"
)

type BigSequenceNumerics struct {
	bigbase.BigBaseNumerics
	area int
}

// Check that BigSequenceNumerics implements SequenceNumerics interface
var _ sequence.SequenceNumerics = (*BigSequenceNumerics)(nil)

func Make(app bigbase.RenderApplication) BigSequenceNumerics {
	w, h := app.PictureDimensions()
	return BigSequenceNumerics{
		BigBaseNumerics: bigbase.Make(app),
		area: int(w * h),
	}
}

func (bsn *BigSequenceNumerics) Sequence() []base.PixelMember {
	imageLeft, imageTop := bsn.PictureMin()
	imageRight, imageBottom := bsn.PictureMax()
	iterlim := bsn.IterateLimit

	area := (imageRight - imageLeft) * (imageBottom - imageTop)
	out := make([]base.PixelMember, area)

	pos := bigbase.BigComplex{
		R: bsn.RealMin,
	}
	count := 0
	for i := imageLeft; i < imageRight; i++ {
		pos.I = bsn.ImagMax
		for j := imageTop; j < imageBottom; j++ {
			member := bigbase.BigMandelbrotMember{
				C: &pos,
				SqrtDivergeLimit: &bsn.SqrtDivergeLimit,
			}
			member.Mandelbrot(iterlim)
			out[count] = base.PixelMember{I: i, J: j, Member: member.MandelbrotMember}

			pos.I.Sub(&pos.I, &bsn.Iunit)
			count++
		}
		pos.R.Add(&pos.R, &bsn.Runit)
	}

	return out
}