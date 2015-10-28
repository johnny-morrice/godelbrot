package draw

import (
	"image"
	"testing"
	"functorama.com/demo/base"
)

func TestDrawPoint(t *testing.T) {
	mockDraw := &MockDrawingContext{
		Pic: image.NewNRGBA(image.ZR),
		Col:  NewRedscalePalette(255),
	}
	pixel := base.PixelMember{1, 2, base.BaseMandelbrot{}}
	DrawPoint(mockDraw, pixel)

	if !(mockDraw.TPicture && mockDraw.TColors) {
		t.Error("Expected method not called on mock drawing context:", mockDraw)
	}
}