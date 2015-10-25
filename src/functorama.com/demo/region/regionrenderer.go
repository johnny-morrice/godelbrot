package libgodelbrot

import (
	"image"
)

type RegionRenderStrategy struct {
	app RenderApplication
}

func NewRegionRenderer(app RenderApplication) *RegionRenderStrategy {
	return &RegionRenderStrategy{app: app}
}

// The RegionRenderStrategy implements RenderNumerics with this method that
// draws the Mandelbrot set uses a "similar rectangles" optimization
func (renderer RegionRenderStrategy) Render() (image.NRGBA, error) {
	// The numerics system is by default a region covering the whole image
	initialRegion := renderer.app.RegionNumerics()
	uniformRegions, smallRegions := SubdivideRegions(initialRegion)

	draw := renderer.app.Draw()

	// Draw uniform regions first
	for _, region := range uniformRegions {
		region.ClaimExtrinsics()
		DrawUniform(draw, region)
	}

	// Add detail from the small regions next
	for _, region := range smallRegions {
		RenderSequentialRegion(region)
	}

	return draw.Picture()
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

		splitee.ClaimExtrinsics()
		// There are three things that can happen to a region...
		//
		// A. The region can be so small that we divide no further
		if Collapse(splitee) {
			smallRegions = append(smallRegions, splitee)
		} else {
			// If the region is not too small, two things can happen
			// B. The region needs subdivided because it covers distinct parts of the plane
			if Subdivide(splitee) {
				splittingRegions = append(splittingRegions, splitee.Children()...)
				// C. The region need not be divided
			} else {
				completeRegions = append(completeRegions, splitee)
			}
		}
	}

	return completeRegions, smallRegions
}
