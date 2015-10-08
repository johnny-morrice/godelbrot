package libgodelbrot

import (
	"image"
)

type SequentialRenderContext interface {
	MandelbrotSequence()
	SetSequencerDrawingContext(drawingContext DrawingContext)
}

func SequentialRender(config SeqentialRenderContext, palette Palette) (*image.NRGBA, error) {
	pic := config.BlankImage()
	config.SequentialRenderImage(config, CreateContext(config, palette, pic))
	return pic, nil
}

func SequentialRenderImage(config SequentialRenderContext, drawingContext DrawingContext) {
	config.SetSequencerDrawingContext(drawingContext)
	config.MandelbrotSequence()
}