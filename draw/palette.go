package draw

import (
	"image/color"
	"github.com/johnny-morrice/godelbrot/base"
)

type Palette interface {
	Color(point base.MandelbrotMember) color.NRGBA
}

type PaletteFactory func(interateLimit uint8) Palette
