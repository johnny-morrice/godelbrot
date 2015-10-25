package sequence

import (
	"image"
	"functorama.com/demo/draw"
)

type SequenceRenderStrategy struct {
	numerics SequenceNumerics	
	context draw.DrawingContext
	iterateLimit uint8
}

func NewSequenceRenderer(app RenderApplication) SequenceRenderStrategy {
	return SequenceRenderStrategy{
		numerics: app.Factory().Build(),
		context: app.DrawingContext(),
		iterateLimit: app.IterateLimit(),
	}
}

// The SequenceRenderStrategy implements RenderContext as it draws the
// Mandelbrot set line by line
func (renderer SequenceRenderStrategy) Render() (*image.NRGBA, error) {
	renderer.numerics.ImageDrawSequencer(renderer.context)
	renderer.numerics.MandelbrotSequence(renderer.iterateLimit)
	return renderer.context.Picture(), nil
}
