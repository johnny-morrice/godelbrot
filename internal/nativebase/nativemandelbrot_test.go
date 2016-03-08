package nativebase

import (
	"testing"
)

func TestMandelbrotSanity(t *testing.T) {
	const origin complex128 = 0
	const non complex128 = 2 + 4i
	const iterateLimit uint8 = 255
	const sqrtDivergeLimit float64 = 2

	originMember := NativeEscapeValue{C: origin, SqrtDivergeLimit: sqrtDivergeLimit}
	nonMember := NativeEscapeValue{C: non, SqrtDivergeLimit: sqrtDivergeLimit}

	originMember.Mandelbrot(iterateLimit)
	nonMember.Mandelbrot(iterateLimit)

	if !originMember.InSet {
		t.Error("Expected origin to be in Mandelbrot set")
	}

	if nonMember.InSet {
		t.Error("Expected ", nonMember, " to be outside Mandelbrot set")
	}

	if nonMember.InvDiv >= iterateLimit {
		t.Error("Expected negativeMembership to have InvDivergence below IterateLimit")
	}

}
