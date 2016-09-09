package draw

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"image/color"
)

type Palette interface {
	Color(point base.EscapeValue) color.NRGBA
}

type PaletteFactory func(interateLimit uint8) Palette
