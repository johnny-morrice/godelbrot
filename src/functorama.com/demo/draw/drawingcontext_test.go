package draw

import (
	"image"
	"testing"
)

func TestDrawUniform(t *testing.T) {
	mockDraw := setupMockContext()
	mockRegion := mockRegionNumerics{}
	DrawUniform(mockDraw, mockRegion)

	if !(mockDraw.TPicture && mockDraw.TColors) {
		t.Error("Expected method not called on mock drawing context:", mockDraw)
	}

	if !(mockRegion.tRegionMember && mockRegion.tRect) {
		t.Error("Expected method not called on mock region numerics:", mockRegion)
	}
}

func TestDrawPoint(t *testing.T) {
	mockDraw := setupMockContext()
	pixel := PixelMember{1, 2, BaseMandelbrotMember{}}
	DrawPoint(mockDraw, pixel)

	if !(mockDraw.TPicture && mockDraw.TColors) {
		t.Error("Expected method not called on mock drawing context:", mockDraw)
	}
}

func setupMockContext() MockDrawingContext {
	return MockDrawingContext{
		Pic: image.NewNRGBA(ZR),
		Col:  NewRedscalePalette(255),
	}
}
