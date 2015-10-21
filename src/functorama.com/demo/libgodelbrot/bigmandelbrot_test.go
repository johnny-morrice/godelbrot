package libgodelbrot

import (
    "testing"
    "math/big"
)

func TestBigMandelbrotSanity(t *testing.T) {
    origin := BigComplex{CreateBigFloat(0.0, Prec64), CreateBigFloat(0.0, Prec64)}
    non := BigComplex{CreateBigFloat(2.0, Prec64), CreateBigFloat(4, Prec64)}
    divergeLimit := CreateBigFloat(4.0, Prec64)
    const iterateLimit uint8 = 255

    originMember := bigMandelbrotHelper(origin)
    nonMember := bigMandelbrotHelper(non)

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