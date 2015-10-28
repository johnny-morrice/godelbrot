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
	limitF := float64(limit)
	linear := float64(limit - index)
	halfLimit := limitF / 2.0
	quatLimit := limitF / 4.0
	qx := linear - halfLimit
	qa := limitF / (halfLimit * halfLimit)
	quadratic := qa * qx * qx
	cx := linear - quatLimit
	ca := -limitF / (quatLimit * quatLimit * quatLimit)
	cubic := ca * cx * cx * cx
	interval := 255.0 / limitF
	return color.NRGBA{
		R: uint8(linear * interval),
		G: uint8(quadratic * interval),
		B: uint8(cubic * interval),
		A: 255,
	}
}
