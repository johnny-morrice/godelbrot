package libgodelbrot

import (
	"image"
	"image/draw"
)

type DrawingContext interface {
	Picture() *image.NRGBA
	Colors() Palette
}

// DrawUniform draws a rectangle of uniform colour on to the image.
func DrawUniform(context DrawingContext, region RegionNumerics) {
	member := region.RegionMember()
	color := context.Colors().Color(member)
	uniform := image.NewUniform(color)
	rect := region.Rect()
	draw.Draw(context.Picture(), rect, uniform, image.ZP, draw.Src)
}

// DrawPoint draws a single point on to the image.
func DrawPoint(context DrawingContext, pixel PixelMember) {
	color := context.Colors().Color(pixel.member)
	context.Picture().Set(pixel.I, pixel.J, color)
}
