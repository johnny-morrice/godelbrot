package libgodelbrot

import (
	"image"
)

type SequentialRenderStrategy struct {
	app *GodelbrotApp
}

func NewSequentialRenderer(app *GodelbrotApp) *RegionRenderStrategy {
	return &SequentialRenderStrategy{app: meditator}
}

// The SequentialRenderStrategy implements RenderContext as it draws the
// Mandelbrot set line by line
func (renderer SequentialRenderStrategy) Render() (image.NRGBA, error) {
	numerics := renderer.app.SequentialNumerics()
	numerics.Initialize()
	numerics.ImageDrawSequencer()
	numerics.MandelbrotSequence()
	return numerics.Picture(), nil
}
