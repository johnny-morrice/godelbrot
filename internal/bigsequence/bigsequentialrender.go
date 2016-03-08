package bigsequence

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/bigbase"
	"github.com/johnny-morrice/godelbrot/internal/sequence"
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
		R: bsn.MakeBigFloat(0.0),
		I: bsn.MakeBigFloat(0.0),
	}
	pos.R.Copy(&bsn.RealMin)
	count := 0
	member := bigbase.BigEscapeValue{
		SqrtDivergeLimit: &bsn.SqrtDivergeLimit,
		Prec: bsn.Precision,
	}
	for i := ileft; i < iright; i++ {
		pos.I.Copy(&bsn.ImagMax)
		for j := itop; j < ibott; j++ {
			member.C = &pos
			member.Mandelbrot(iterlim)
			out[count] = base.PixelMember{I: i, J: j, Member: member.EscapeValue}

			pos.I.Sub(&pos.I, &bsn.Iunit)
			count++
		}
		pos.R.Add(&pos.R, &bsn.Runit)
	}

	return out
}