package libgodelbrot

import (
    "image/color"
)


type RedscalePalette struct {
    CachePalette
}

func NewRedscalePalette(iterateLimit uint8) RedscalePalette {
    black := color.NRGBA{
        R: 0, G: 0, B: 0, A: 0,
    }
    return RedscalePalette{
        NewCachePalette(iterateLimit, black, redscaleCacher),
    }
}

// Cache redscale colour values
func redscaleCacher(limit uint8, index uint8) color.NRGBA {
    return color.NRGBA{
        R: index * (255 / limit),
        G: 0,
        B: 0,
        A: 255,
    }
}