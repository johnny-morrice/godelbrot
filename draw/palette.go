package draw

import (
	"image/color"
	"functorama.com/demo/base"
)

type Palette interface {
	Color(point base.MandelbrotMember) color.NRGBA
}

type PaletteFactory func(interateLimit uint8) Palette
