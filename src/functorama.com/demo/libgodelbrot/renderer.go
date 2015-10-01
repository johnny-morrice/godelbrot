package libgodelbrot

import (
	"image"
)

type Renderer func(config *RenderConfig, palette Palette) (*image.NRGBA, error)
