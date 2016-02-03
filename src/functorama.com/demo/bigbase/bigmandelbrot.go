package bigbase

import (
	"math/big"
	"functorama.com/demo/base"
)

type BigMandelbrotMember struct {
	base.MandelbrotMember
	C                *BigComplex
	SqrtDivergeLimit *big.Float
	Prec uint
}

func (member *BigMandelbrotMember) Mandelbrot(iterateLimit uint8) {
	z := MakeBigComplex(0.0, 0.0, member.Prec)
	aa := MakeBigFloat(0.0, member.Prec)
	bb := MakeBigFloat(0.0, member.Prec)
	ab := MakeBigFloat(0.0, member.Prec)
	i := uint8(0)
	for ; i < iterateLimit && withinMandLimit(&z, member.SqrtDivergeLimit); i++ {
		aa.Mul(z.Real(), z.Real())

		bb.Mul(z.Imag(), z.Imag())
		ab.Mul(z.Real(), z.Imag())

		z.R.Copy(aa.Sub(&aa, &bb))
		z.I.Copy(ab.Add(&ab, &ab))

		z.Add(&z, member.C)
	}

	member.InSet = i >= iterateLimit
	member.InvDiv = i
}

func withinMandLimit(z *BigComplex, limit *big.Float) bool {
	// Approximate cmplx.Abs
	negLimit := MakeBigFloat(0.0, limit.Prec())
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

func nativec(c BigComplex) complex128 {
	fr, _ := c.Real().Float64()
	fi, _ := c.Imag().Float64()
	return complex(fr, fi)
}
