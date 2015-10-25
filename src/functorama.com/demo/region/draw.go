package region

import(
	image "image/draw"
	"functorama.com/demo/draw"
)

// DrawUniform draws a rectangle of uniform colour on to the image.
func DrawUniform(context draw.DrawingContext, region RegionNumerics) {
	member := region.RegionMember()
	color := context.Colors().Color(member)
	uniform := image.NewUniform(color)
	rect := region.Rect()
	image.Draw(context.Picture(), rect, uniform, image.ZP, draw.Src)
}