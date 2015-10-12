package libgodelbrot

import (
	"image"
	"image/draw"
)

type DrawingContext interface {
    Picture() *image.NRGBA
	Paint() Palette
}

type RegionDrawingContext interface {
	DrawingContext
	RegionMember() MandelbrotMember
	Rect() image.Rectangle
}

func (context RegionDrawingContext) DrawUniform() {
	member := context.RegionMember()
	color := context.Paint().Color(member)
	uniform := image.NewUniform(color)
	rect := context.Rect()
	draw.Draw(context.Picture(), rect, uniform, image.ZP, draw.Src)
}

func (context DrawingContext) DrawPointAt(i int, j int, member MandelbrotMember) {
	color := context.Paint().Color(member)
	context.Picture().Set(i, j, color)
}
