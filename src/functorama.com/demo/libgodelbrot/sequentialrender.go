package libgodelbrot

import (
	"image"
)

type SequentialRenderStrategy struct {
    Meditator *ContextMediator
}

func NewSequentialRenderer(meditator *ContextMediator) *RegionRenderStrategy {
    return &SequentialRenderStrategy{Meditator: meditator}
}

// The SequentialRenderStrategy implements RenderContext as it draws the
// Mandelbrot set line by line
func (renderer SequentialRenderStrategy) Render() (image.NRGBA, error) {
    numerics := renderer.Meditator.SequentialNumerics()
    numerics.ImageDrawSequencer()
    numerics.MandelbrotSequence()
    return numerics.Picture(), nil
}