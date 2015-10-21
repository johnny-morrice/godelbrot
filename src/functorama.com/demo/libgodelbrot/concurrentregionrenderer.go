package libgodelbrot

import (
	"fmt"
	"image"
)

func NewConcurrentRegionRenderer(app *GodelbrotApplication) *RenderTracker {
	return NewRenderTracker(app)
}
