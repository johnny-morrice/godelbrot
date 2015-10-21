package libgodelbrot

import (
	"image"
	"image/draw"
)

type DrawingContext interface {
    Picture() *image.NRGBA
	Colors() Palette
}

type RegionDrawingContext interface {
	DrawingContext
	RegionMember() MandelbrotMember
	Rect() image.Rectangle
}

func (context RegionDrawingContext) DrawUniform() {
	member := context.RegionMember()
	color := context.Colors().Color(member)
	uniform := image.NewUniform(color)
	rect := context.Rect()
	draw.Draw(context.Picture(), rect, uniform, image.ZP, draw.Src)
}

func (context DrawingContext) DrawPointAt(i int, j int, member MandelbrotMember) {
	color := context.Colors().Color(member)
	context.Picture().Set(i, j, color)
}
