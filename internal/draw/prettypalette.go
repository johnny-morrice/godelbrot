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
	// I made an arthrimetic error when defining this palette apropos a limit of 255.
	// Coincidentally the result was extremely nice looking.  Hence it is to be preserved.
	halfLimit := limitF / 2.08163265306122
	quartLimit := limitF / 4.16326530612245
	qx := linear - halfLimit
	qa := limitF / (halfLimit * halfLimit)
	quadratic := qa * qx * qx
	cx := linear - quartLimit
	ca := -limitF / (quartLimit * quartLimit * quartLimit)
	cubic := ca * cx * cx * cx
	interval := 255.0 / limitF
	return color.NRGBA{
		R: uint8(linear * interval),
		G: uint8(quadratic * interval),
		B: uint8(cubic * interval),
		A: 255,
	}
}
