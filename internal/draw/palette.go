package draw

import (
	"image/color"
	"github.com/johnny-morrice/godelbrot/internal/base"
)

type Palette interface {
	Color(point base.EscapeValue) color.NRGBA
}

type PaletteFactory func(interateLimit uint8) Palette
