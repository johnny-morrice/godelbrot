package region

import (
	"github.com/johnny-morrice/godelbrot/internal/draw"
	"image"
	"testing"
)

func TestDrawUniform(t *testing.T) {
	mockDraw := &draw.MockDrawingContext{
		Pic: image.NewNRGBA(image.ZR),
		Col: draw.NewRedscalePalette(255),
	}
	mockRegion := &MockNumerics{}
	DrawUniform(mockDraw, mockRegion)

	if !(mockDraw.TPicture && mockDraw.TColors) {
		t.Error("Expected method not called on mock drawing context:", mockDraw)
	}

	if !(mockRegion.TRegionMember && mockRegion.TRect) {
		t.Error("Expected method not called on mock region numerics:", mockRegion)
	}
}
