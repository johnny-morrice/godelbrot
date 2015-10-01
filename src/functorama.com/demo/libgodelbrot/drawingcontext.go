package libgodelbrot

import (
	"image"
	"image/draw"
)

type DrawingContext struct {
	Pic          *image.NRGBA
	ColorPalette Palette
	Config       *RenderConfig
}

func CreateContext(config *RenderConfig, palette Palette, pic *image.NRGBA) DrawingContext {
	return DrawingContext{
		Pic:          pic,
		ColorPalette: palette,
		Config:       config,
	}
}

func (context DrawingContext) DrawUniform(region Region) {
	member := region.midPoint.membership
	color := context.ColorPalette.Color(member)
	uniform := image.NewUniform(color)
	rect := region.Rect(context.Config)
	draw.Draw(context.Pic, rect, uniform, image.ZP, draw.Src)
}

func (context DrawingContext) DrawPointAt(i int, j int, member MandelbrotMember) {
	color := context.ColorPalette.Color(member)
	context.Pic.Set(i, j, color)
}
