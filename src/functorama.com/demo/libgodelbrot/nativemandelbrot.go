package libgodelbrot

import (
	"math"
)

type NativeMandelbrotMember struct {
    BaseMandelbrot
    C complex128
}

func (member *NativeMandelbrotMember) Mandelbrot(iterateLimit uint8, divergeLimit float64) MandelbrotMember {
    var z complex128 = 0
    sqrtDl := math.Sqrt(divergeLimit)
    c := member.C
    i := uint8(0)
    for ; i < iterateLimit && z.WithinMandLimit(sqrtDl); i++ {
        z = (z * z) + c
    }

    member.InSet = i >= iterateLimit,
    member.InvDivergence = i
}

func (z complex128) WithinMandLimit(limit float64) bool {
    // Approximate cmplx.Abs
    negLimit := -limit
    x := real(z)
    y := imag(z)
    return x < limit && x > negLimit && y < limit && y > negLimit
}