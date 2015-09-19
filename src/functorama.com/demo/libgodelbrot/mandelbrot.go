package libgodelbrot

import (
    "math/cmplx"
)

const limit int = 2

type MandelbrotMember struct {
    InSet bool
    Divergence uint8
    C complex128
}

func Mandelbrot(c complex128) MandebrotMember {
    var z complex128 = 0
    i := 0
    for ; i < 255 && cmplx.Abs(z) < 4; i++ {
        z = (z * z) + c
    }

    return MandelbrotMember{
        InSet: i >= 255
        Divergence: i
        C: c
    }
}