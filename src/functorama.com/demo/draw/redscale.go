package draw

import (
	"image/color"
)

type RedscalePalette struct {
	CachePalette
}

func NewRedscalePalette(iterateLimit uint8) Palette {
	black := color.NRGBA{
		R: 0, G: 0, B: 0, A: 255,
	}
	return RedscalePalette{
		NewCachePalette(iterateLimit, black, redscaleCacher),
	}
}

// Cache redscale colour values
func redscaleCacher(limit uint8, index uint8) color.NRGBA {
	calibIndex := float64(limit - index)
	interval := 255.0 / float64(limit)
	return color.NRGBA{
		R: uint8(calibIndex * interval),
		G: 0,
		B: 0,
		A: 255,
	}
}
