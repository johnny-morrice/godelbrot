package bigbase

import (
	"math/big"
	"functorama.com/demo/base"
)

type BigMandelbrotMember struct {
	base.BaseMandelbrot
	C                *BigComplex
	SqrtDivergeLimit *big.Float
	Prec uint
}

var _ base.MandelbrotMember = (*BigMandelbrotMember)(nil)

func (member *BigMandelbrotMember) Mandelbrot(iterateLimit uint8) {
	z := CreateBigComplex(0.0, 0.0, member.Prec)
	aa := MakeBigFloat(0.0, member.Prec)
	bb := MakeBigFloat(0.0, member.Prec)
	ab := MakeBigFloat(0.0, member.Prec)
	i := uint8(0)
	for ; i < iterateLimit && withinMandLimit(&z, member.SqrtDivergeLimit); i++ {
		aa.Set(z.Real())
		aa.Mul(&aa, &aa)

		bb.Set(z.Imag())
		bb.Mul(&bb, &bb)

		ab.Set(z.Real())
		ab.Mul(&ab, z.Imag())

		z.R = *aa.Sub(&aa, &bb)
		z.I = *ab.Add(&ab, &ab)

		z.Add(&z, member.C)
	}

	member.InSet = i >= iterateLimit
	member.InvDivergence = i
}

func withinMandLimit(z *BigComplex, limit *big.Float) bool {
	// Approximate cmplx.Abs
	negLimit := big.Float{}
	negLimit.SetPrec(limit.Prec())
	negLimit.Neg(limit)

	r := z.Real()
	i := z.Imag()

	rLimCmp := r.Cmp(limit)
	rNegLimCmp := r.Cmp(&negLimit)
	iLimCmp := i.Cmp(limit)
	iNegLimCmp := i.Cmp(&negLimit)

	within := rLimCmp == -1 && rNegLimCmp == 1
	within = within && iLimCmp == -1 && iNegLimCmp == 1
	return within
}
