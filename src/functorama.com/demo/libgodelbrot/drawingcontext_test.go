package libgodelbrot

import (
	"image"
	"testing"
)

type mockDrawingContext struct {
	tPicture bool
	tColors  bool

	picture *image.NRGBA
	colors  Palette
}

func (mock *mockDrawingContext) Picture() *image.NRGBA {
	mock.tPicture = true
	return mock.picture
}

func (mock *mockDrawingContext) Colors() Palette {
	mock.tColors = true
	return mock.colors
}

func TestDrawUniform(t *testing.T) {
	mockDraw := setupMockContext()
	mockRegion := mockRegionNumerics{}
	DrawUniform(mockDraw, mockRegion)

	if !(mockDraw.tPicture && mockDraw.tColors) {
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

	if !(mockDraw.tPicture && mockDraw.tColors) {
		t.Error("Expected method not called on mock drawing context:", mockDraw)
	}
}

func setupMockContext() mockDrawingContext {
	return mockDrawingContext{
		picture: image.NewNRGBA(ZR),
		colors:  NewRedscalePalette(255),
	}
}
