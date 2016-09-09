package base

import (
	"image"
)

// Numerics that are aware of a picture and of the Mandelbrot iteration limit
type BaseNumerics struct {
	WholeWidth  uint
	WholeHeight uint
	PicXMin     int
	PicXMax     int // exclusive maximum
	PicYMin     int
	PicYMax     int // exclusive maximum
}

func Make(app RenderApplication) BaseNumerics {
	w, h := app.PictureDimensions()
	base := BaseNumerics{}
	base.WholeWidth = w
	base.WholeHeight = h
	base.RestorePicBounds()
	return base
}

func (base *BaseNumerics) PictureMin() (int, int) {
	return base.PicXMin, base.PicYMin
}

func (base *BaseNumerics) PictureMax() (int, int) {
	return base.PicXMax, base.PicYMax
}

func (base *BaseNumerics) RestorePicBounds() {
	base.setWidth(base.WholeWidth)
	base.setHeight(base.WholeHeight)
}

func (base *BaseNumerics) setWidth(width uint) {
	base.PicXMax = int(width)
	base.PicXMin = 0
}

func (base *BaseNumerics) setHeight(height uint) {
	base.PicYMax = int(height)
	base.PicYMin = 0
}

// Change the drawing context to a sub-part of the image
func (base *BaseNumerics) PictureSubImage(rect image.Rectangle) {
	base.PicXMin = rect.Min.X
	base.PicYMin = rect.Min.Y
	base.PicXMax = rect.Max.X
	base.PicYMax = rect.Max.Y
}
