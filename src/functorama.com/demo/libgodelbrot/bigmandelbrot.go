package libgodelbrot

import (
    "math/big"
)

type BigMandelbrotMember struct {
    BaseMandelbrot
    C BigComplex
}

func (member *BigMandelbrotMember) Mandelbrot(iterateLimit uint8, divergeLimit float64) {
    prec := member.C.Prec()
    z := NewBigComplexNative(0.0, 0.0, prec)
    sqrtDl := math.Sqrt(divergeLimit)
    aa := NewFloatPrec(0, prec)
    bb := NewFloatPrec(0, prec)
    ab := NewFloatPrec(0, prec)
    i := uint8(0)
    c := member.C
    for ; i < iterateLimit && z.WithinMandLimit(sqrtDl); i++ {
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

    member.InSet = i >= iterateLimit,
    member.InvDivergence = i
}

func (z BigComplex) WithinMandLimit(limit float64) bool {
    // Approximate cmplx.Abs
    negLimit := -limit
    x := z.R
    y := z.Y
    return x.Lt(limit) && x.Gt(negLimit) && y.Lt(limit) && y.Gt(negLimit)
}