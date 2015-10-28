package draw

import (
    "image"
)

type MockDrawingContext struct {
    TPicture bool
    TColors  bool

    Pic *image.NRGBA
    Col  Palette
}

func (mock *MockDrawingContext) Picture() *image.NRGBA {
    mock.TPicture = true
    return mock.Pic
}

func (mock *MockDrawingContext) Colors() Palette {
    mock.TColors = true
    return mock.Col
}