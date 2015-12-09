package bigbase

import (
	"testing"
)

func TestBigMandelbrotSanity(t *testing.T) {
	origin := BigComplex{MakeBigFloat(0.0, testPrec), MakeBigFloat(0.0, testPrec)}
	non := BigComplex{MakeBigFloat(2.0, testPrec), MakeBigFloat(4, testPrec)}
	sqrtDL := MakeBigFloat(2.0, testPrec)
	const iterateLimit uint8 = 255

	originMember := BigMandelbrotMember{
		C: &origin,
		SqrtDivergeLimit: &sqrtDL,
		Prec: testPrec,
	}
	nonMember := BigMandelbrotMember{
		C: &non,
		SqrtDivergeLimit: &sqrtDL,
		Prec: testPrec,
	}

	originMember.Mandelbrot(iterateLimit)
	nonMember.Mandelbrot(iterateLimit)

	if !originMember.SetMember() {
		t.Error("Expected origin to be in Mandelbrot set")
	}

	if nonMember.SetMember() {
		t.Error("Expected ", nonMember, " to be outside Mandelbrot set")
	}

	if nonMember.InverseDivergence() >= iterateLimit {
		t.Error("Expected negativeMembership to have InvDivergence below IterateLimit")
	}
}
