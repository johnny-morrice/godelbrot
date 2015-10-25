package base

import (
	"image"
)

// A reusable notion of collapsable regions
type BaseRegionNumerics struct {
	GlitchSamples int
	Collapse      int
}

func (collapse BaseRegionNumerics) CollapseSize() int {
	return collapse.Collapse
}

// Numerics that are aware of a picture and of the Mandelbrot iteration limit
type BaseNumerics struct {
	PicXMin int
	PicXMax int // exclusive maximum
	PicYMin int
	PicYMax int // exclusive maximum
}

func CreateBaseNumerics(app RenderApplication) BaseNumerics {
	width, height := app.PictureDimensions()
	return BaseNumerics{
		PicXMin: 0,
		PicXMax: int(width),
		PicYMin: 0,
		PicYMax: int(height),
	}
}

func (base *BaseNumerics) PictureMin() (int, int) {
	return base.PicXMin, base.PicYMin
}

func (base *BaseNumerics) PictureMax() (int, int) {
	return base.PicXMax, base.PicYMax
}

// Change the drawing context to a sub-part of the image
func (native *BaseNumerics) SubImage(rect image.Rectangle) {
	native.PicXMin = rect.Min.X
	native.PicYMin = rect.Min.Y
	native.PicXMax = rect.Max.X
	native.PicYMax = rect.Max.Y
}