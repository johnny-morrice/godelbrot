package draw

import (
	"image"
	"functorama.com/demo/base"
)

type DrawingContext interface {
	Picture() *image.NRGBA
	Colors() Palette
}

// DrawPoint draws a single point on to the image.
func DrawPoint(context DrawingContext, pixel base.PixelMember) {
	color := context.Colors().Color(pixel.Member)
	context.Picture().Set(pixel.I, pixel.J, color)
}
