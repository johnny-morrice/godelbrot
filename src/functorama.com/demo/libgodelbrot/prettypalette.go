package libgodelbrot

import (
    "image/color"
)

type PrettyPalette struct {
    CachePalette
}

func NewPrettyPalette(iterateLimit uint8) Palette {
    black := color.NRGBA{
        R: 0, G: 0, B: 0, A: 255,
    }
    return PrettyPalette{
        NewCachePalette(iterateLimit, black, prettyCacher),
    }
}

// Cache redscale colour values
func prettyCacher(limit uint8, index uint8) color.NRGBA {
    limitF := float64(limit)
    linear := limitF - float64(index)
    mid := float64(limit) / 2.0
    quart := mid / 2.0
    qx := linear - mid
    qa := limitF / (mid * mid)
    quadratic := qa * qx * qx
    cx := linear - quart
    ca := - limitF / (quart * quart * quart)
    cubic := ca * cx * cx * cx
    interval := 255.0 / limitF
    return color.NRGBA{
        R: uint8(linear * interval),
        G: uint8(quadratic * interval),
        B: uint8(cubic * interval),
        A: 255,
    }
}