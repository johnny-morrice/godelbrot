package draw

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
	linear := float64(limit - index)
	qx := linear - 122.5
	qa := 255.0 / (122.5 * 122.5)
	quadratic := qa * qx * qx
	cx := linear - 61.25
	ca := -255 / (61.25 * 61.25 * 61.25)
	cubic := ca * cx * cx * cx
	interval := 255.0 / float64(limit)
	return color.NRGBA{
		R: uint8(linear * interval),
		G: uint8(quadratic * interval),
		B: uint8(cubic * interval),
		A: 255,
	}
}
