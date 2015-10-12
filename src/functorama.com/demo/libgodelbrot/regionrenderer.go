package libgodelbrot

import (
	"image"
)

type RegionRenderStrategy struct {
	Context *ContextFacade
}

func NewRegionRenderer(meditator *ContextFacade) *RegionRenderStrategy {
	return &RegionRenderStrategy{Context: meditator}
}

// The RegionRenderStrategy implements RenderNumerics with this method that
// draws the Mandelbrot set uses a "similar rectangles" optimization
func (renderer RegionRenderStrategy) Render() (image.NRGBA, error) {
	numerics := renderer.Context.RegionNumerics()
	initialRegion := numerics.WholeRegion()
	uniformRegions, smallRegions := SubdivideRegions(initialRegion)

	// Draw uniform regions first
	for _, region := range uniformRegions {
		drawingNumerics.DrawUniform(region)
	}

	// Add detail from the small regions next
	for _, region := range smallRegions {
		RenderSequentialRegion(region)
	}

	return numerics.Picture()
}

func SubdivideRegions(whole RegionNumerics) ([]RegionNumerics, []RegionNumerics) {
	// Lots of preallocated space for regions and region pointers
	completeRegions := make([]RegionNumerics, 0, allocMedium)
	smallRegions := make([]RegionNumerics, 0, allocMedium)
	splittingRegions := make([]RegionNumerics, 1, allocMedium)

	// Split regions
	splittingRegions[0] = whole
	for i := 0; i < len(splittingRegions); i++ {
		splitee := splittingRegions[i]

		// There are three things that can happen to a region...
		//
		// A. The region can be so small that we divide no further
		if splitee.Collapse(config) {
			smallRegions = append(smallRegions, splitee)
		} else {
			// If the region is not too small, two things can happen
			// B. The region needs subdivided because it covers distinct parts of the plane
			if Subdivide(splitee) {
				splittingRegions = append(splittingRegions, splitee.Children() ...)
				// C. The region need not be divided
			} else {
				completeRegions = append(completeRegions, splitee)
			}
		}
	}

	return completeRegions, smallRegions
}
