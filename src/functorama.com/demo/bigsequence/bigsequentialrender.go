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
	ileft, itop := bsn.PictureMin()
	iright, ibott := bsn.PictureMax()
	iterlim := bsn.IterateLimit

	area := (iright - ileft) * (ibott - itop)
	out := make([]base.PixelMember, area)

	pos := bigbase.BigComplex{
		R: bsn.RealMin,
	}
	count := 0
	member := bigbase.BigMandelbrotMember{
		SqrtDivergeLimit: &bsn.SqrtDivergeLimit,
		Prec: bsn.Precision,
	}
	for i := ileft; i < iright; i++ {
		pos.I = bsn.ImagMax
		for j := itop; j < ibott; j++ {
			member.C = &pos
			member.Mandelbrot(iterlim)
			out[count] = base.PixelMember{I: i, J: j, Member: member.MandelbrotMember}

			pos.I.Sub(&pos.I, &bsn.Iunit)
			count++
		}
		pos.R.Add(&pos.R, &bsn.Runit)
	}

	return out
}