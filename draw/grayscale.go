package draw

import (
    "image/color"
)

type GrayscalePalette struct {
    CachePalette
}

func NewGrayscalePalette(iterateLimit uint8) Palette {
    white := color.NRGBA{
        R: 255, G: 255, B: 255, A: 255,
    }
    return GrayscalePalette{
        NewCachePalette(iterateLimit, white, grayCache),
    }
}

// Cache redscale colour values
func grayCache(limit uint8, index uint8) color.NRGBA {
    calibIndex := float64(index)
    interval := 255.0 / float64(limit)
    gray := uint8(calibIndex * interval)
    return color.NRGBA{
        R: gray,
        G: gray,
        B: gray,
        A: 255,
    }
}
