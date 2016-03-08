package sequence

import (
	"image"
	"github.com/johnny-morrice/godelbrot/internal/draw"
)

type SequenceRenderStrategy struct {
	numerics SequenceNumerics
	context draw.DrawingContext
}

func Make(app RenderApplication) SequenceRenderStrategy {
	return SequenceRenderStrategy{
		numerics: app.SequenceNumericsFactory().Build(),
		context: app.DrawingContext(),
	}
}

// The SequenceRenderStrategy implements RenderContext as it draws the
// Mandelbrot set line by line
func (srs SequenceRenderStrategy) Render() (*image.NRGBA, error) {
	ImageSequence(srs.numerics, srs.context)
	return srs.context.Picture(), nil
}
