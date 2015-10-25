package nativebase

import (
	"testing"
)

func TestMandelbrotSanity(t *testing.T) {
	const originMember complex128 = 0
	const nonMember complex128 = 2 + 4i
	const iterateLimit uint8 = 255
	const divergeLimit float64 = 4.0

	originMember := CreateNativeMandelbrotMember(origin)
	nonMember := CreateNativeMandelbrotMember(non)

	originMember.Mandelbrot(iterateLimit)
	nonMember.Mandelbrot(iterateLimit)

	if !originMember.InSet() {
		t.Error("Expected origin to be in Mandelbrot set")
	}

	if non.InSet() {
		t.Error("Expected ", nonMember, " to be outside Mandelbrot set")
	}

	if non.InverseDivergence() >= iterateLimit {
		t.Error("Expected negativeMembership to have InvDivergence below IterateLimit")
	}

}
