package libgodelbrot

import (
    "math/cmplx"
    "image/color"
)


type MandelbrotMember struct {
    InSet bool
    InvDivergence uint8
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
        InvDivergence: i,
        C: c,
    }
}

type MandelbrotPalette []color.NRGBA

func NewRedscalePalette() MandelbrotPalette {
    palette := make([]color.NRGBA, 256, 256)
    for i := 0; i < 255; i++ {
        palette[i] = redscaleColor(uint8(i))
    }
    return palette
}

func (palette MandelbrotPalette) Lookup(member MandelbrotMember) color.NRGBA {
    if member.InSet {
        return color.NRGBA{
            R: 0,
            G: 0,
            B: 0,
            A: 255,
        }
    } else {
        return palette[member.InvDivergence]
    }
}

func redscaleColor(index uint8) color.NRGBA {
    return color.NRGBA{
        R: 255 - index,
        G: 0,
        B: 0,
        A: 255,
    }
}