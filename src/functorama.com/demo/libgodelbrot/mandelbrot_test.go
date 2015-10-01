package libgodelbrot

import (
	"testing"
)

func TestMandelbrotSanity(t *testing.T) {
	const originMember complex128 = 0
	const nonMember complex128 = 2 + 4i
	const iterateLimit uint8 = 255
	const divergeLimit float64 = 4.0

	positiveMembership := Mandelbrot(originMember, iterateLimit, divergeLimit)
	negativeMembership := Mandelbrot(nonMember, iterateLimit, divergeLimit)

	if !positiveMembership.InSet {
		t.Error("Expected origin to be in Mandelbrot set")
	}

	if positiveMembership.C != originMember {
		t.Error("Expected Mandelbrot return to contain origin co-ordinate, got ", positiveMembership.C)
	}

	if negativeMembership.InSet {
		t.Error("Expected ", nonMember, " to be outside Mandelbrot set")
	}

	if negativeMembership.C != nonMember {
		t.Error("Expected Mandelbrot return to contain ", nonMember, " got ", negativeMembership.C)
	}

	if negativeMembership.InvDivergence >= iterateLimit {
		t.Error("Expected negativeMembership to have InvDivergence below IterateLimit")
	}

}
