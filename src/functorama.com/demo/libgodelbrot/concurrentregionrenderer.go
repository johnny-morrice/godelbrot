package libgodelbrot

import (
	"fmt"
	"image"
)

func NewConcurrentRegionRenderer(app *GodelbrotApplication) {
	tracker := NewRenderTracker(app)
	return tracker.Render()
}
