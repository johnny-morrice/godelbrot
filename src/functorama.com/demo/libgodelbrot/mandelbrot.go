package libgodelbrot

import (
    "math/cmplx"
)


type MandelbrotMember struct {
    InSet bool
    Divergence uint8
    C complex128
}

func Mandelbrot(c complex128, iterateLimit uint8, divergeLimit float64) MandelbrotMember {
    var z complex128 = 0
    var i uint8 = 0
    for ; i < iterateLimit && cmplx.Abs(z) < divergeLimit; i++ {
        z = (z * z) + c
    }

    return MandelbrotMember{
        InSet: i >= iterateLimit,
        Divergence: i,
        C: c,
    }
}