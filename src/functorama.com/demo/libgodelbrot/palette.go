package libgodelbrot

import (
	"image/color"
)

type Palette interface {
	Color(point MandelbrotMember) color.NRGBA
}

type PaletteFactory func(interateLimit uint8) Palette