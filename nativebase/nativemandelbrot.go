package nativebase

import (
	"functorama.com/demo/base"
)

type NativeMandelbrotMember struct {
	base.MandelbrotMember
	C complex128
	SqrtDivergeLimit float64
}

func (member *NativeMandelbrotMember) Mandelbrot(iterateLimit uint8) {
	var z complex128 = 0
	sqrtDl := member.SqrtDivergeLimit
	c := member.C
	i := uint8(0)
	for ; i < iterateLimit && withinMandLimit(z, sqrtDl); i++ {
		z = (z * z) + c
	}

	member.InSet = i >= iterateLimit
	member.InvDiv = i
}

func withinMandLimit(z complex128, limit float64) bool {
	// Approximate cmplx.Abs
	negLimit := -limit
	x := real(z)
	y := imag(z)
	return x < limit && x > negLimit && y < limit && y > negLimit
}