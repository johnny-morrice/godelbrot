package region

import (
	paint "github.com/johnny-morrice/godelbrot/internal/draw"
	"image"
	"image/draw"
)

// DrawUniform draws a rectangle of uniform colour on to the image.
func DrawUniform(context paint.DrawingContext, region RegionNumerics) {
	member := region.RegionMember()
	color := context.Colors().Color(member)
	uniform := image.NewUniform(color)
	rect := region.Rect()

	draw.Draw(context.Picture(), rect, uniform, image.ZP, draw.Src)
}
