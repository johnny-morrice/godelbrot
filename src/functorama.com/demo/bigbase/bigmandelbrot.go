package bigbase

import (
	"math/big"
)

type BigMandelbrotMember struct {
	BaseMandelbrot
	C                BigComplex
	SqrtDivergeLimit big.Float
}

func CreateBigMandelbrotMember(real, imag) {
	return BigMandelbrotMember{
		C: BigComplex{R: real, I: imag},
	}
}

func (member *BigMandelbrotMember) Mandelbrot(iterateLimit uint8) {
	prec := member.C.Prec()
	z := NewBigComplex(0.0, 0.0, prec)
	aa := NewBigFloat(0.0, prec)
	bb := NewBigFloat(0.0, prec)
	ab := NewBigFloat(0.0, prec)
	i := uint8(0)
	c := member.C
	for ; i < iterateLimit && withinMandLimit(member.SqrtDLimit); i++ {
		aa.Set(z.R)
		aa.Mul(aa, aa)

		bb.Set(z.I)
		bb.Mul(bb, bb)

		ab.Set(z.R)
		ab.Mul(ab, z.I)

		z.R = aa.Sub(aa, bb)
		z.I = ab.Add(ab, ab)

		z.Add(z, c)
	}

	member.InSet = i >= iterateLimit
	member.InvDivergence = i
}

func withinBigMandLimit(z BigComplex, limit big.Float) bool {
	// Approximate cmplx.Abs
	negLimit := -limit
	r := z.R
	i := z.Y

	rLimCmp := r.Cmp(limit)
	rNegLimCmp := r.Cmp(negLimit)
	iLimCmp := i.Cmp(limit)
	iNegLimCmp := i.Cmp(negLimit)

	within := rLimCmp == -1 && rNegLimCmp == 1
	within = within && iLimCmp == -1 && iNegLimCmp == 1
	return within
}
