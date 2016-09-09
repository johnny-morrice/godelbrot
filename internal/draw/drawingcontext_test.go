package draw

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"image"
	"testing"
)

func TestDrawPoint(t *testing.T) {
	mockDraw := &MockDrawingContext{
		Pic: image.NewNRGBA(image.ZR),
		Col: NewRedscalePalette(255),
	}
	pixel := base.PixelMember{1, 2, base.EscapeValue{}}
	DrawPoint(mockDraw, pixel)

	if !(mockDraw.TPicture && mockDraw.TColors) {
		t.Error("Expected method not called on mock drawing context:", mockDraw)
	}
}
