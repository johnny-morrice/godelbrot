package draw

import (
	"image"
	"testing"
	"github.com/johnny-morrice/godelbrot/base"
)

func TestDrawPoint(t *testing.T) {
	mockDraw := &MockDrawingContext{
		Pic: image.NewNRGBA(image.ZR),
		Col:  NewRedscalePalette(255),
	}
	pixel := base.PixelMember{1, 2, base.MandelbrotMember{}}
	DrawPoint(mockDraw, pixel)

	if !(mockDraw.TPicture && mockDraw.TColors) {
		t.Error("Expected method not called on mock drawing context:", mockDraw)
	}
}