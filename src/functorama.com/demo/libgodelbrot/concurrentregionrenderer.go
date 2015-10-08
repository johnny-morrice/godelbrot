package libgodelbrot

import (
	"fmt"
	"image"
)

func ConcurrentRegionRender(config *RenderConfig, palette Palette) (*image.NRGBA, error) {
	pic := config.BlankImage()
	ConcurrentRegionRenderImage(CreateContext(config, palette, pic))
	return pic, nil
}

func ConcurrentRegionRenderImage(drawingContext DrawingContext) {
	tracker := NewRenderTracker(drawingContext)
	tracker.Render()
}
