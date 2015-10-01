package libgodelbrot

import (
	"math"
)

type MandelbrotMember struct {
	InSet         bool
	InvDivergence uint8
	C             complex128
}

func Mandelbrot(c complex128, iterateLimit uint8, divergeLimit float64) MandelbrotMember {
	var z complex128 = 0
	var i uint8 = 0
	sqrtDl := math.Sqrt(divergeLimit)
	for ; i < iterateLimit && withinMandLimit(z, sqrtDl); i++ {
		z = (z * z) + c
	}

	return MandelbrotMember{
		InSet:         i >= iterateLimit,
		InvDivergence: i,
		C:             c,
	}
}

func withinMandLimit(z complex128, limit float64) bool {
	// Approximate cmplx.Abs
	negLimit := -limit
	x := real(z)
	y := imag(z)
	return x < limit && x > negLimit && y < limit && y > negLimit
}
