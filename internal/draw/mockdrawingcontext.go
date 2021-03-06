package draw

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"image"
	"image/color"
)

type MockDrawingContext struct {
	TPicture bool
	TColors  bool

	Pic *image.NRGBA
	Col Palette
}

var _ DrawingContext = (*MockDrawingContext)(nil)

func NewMockDrawingContext(iterateLimit uint8) *MockDrawingContext {
	return &MockDrawingContext{
		Pic: image.NewNRGBA(image.ZR),
		Col: NewRedscalePalette(iterateLimit),
	}
}

func (mock *MockDrawingContext) Picture() *image.NRGBA {
	mock.TPicture = true
	return mock.Pic
}

func (mock *MockDrawingContext) Colors() Palette {
	mock.TColors = true
	return mock.Col
}

type MockPalette struct {
	TColor bool
	Col    color.NRGBA
}

func (mock *MockPalette) Color(point base.EscapeValue) color.NRGBA {
	mock.TColor = true
	return mock.Col
}
