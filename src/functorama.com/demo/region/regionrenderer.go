package region

import (
	"image"
	"functorama.com/demo/base"
)

type RegionRenderStrategy struct {
	app RenderApplication
}

func NewRegionRenderer(app RenderApplication) *RegionRenderStrategy {
	return &RegionRenderStrategy{app: app}
}

// The RegionRenderStrategy implements RenderNumerics with this method that
// draws the Mandelbrot set uses a "similar rectangles" optimization
func (renderer RegionRenderStrategy) Render() (*image.NRGBA, error) {
	// The numerics system is by default a region covering the whole image
	initialRegion := renderer.app.Factory().Build()
	uniformRegions, smallRegions := renderer.SubdivideRegions(initialRegion)

	draw := renderer.app.DrawingContext()

	// Draw uniform regions first
	for _, region := range uniformRegions {
		region.ClaimExtrinsics()
		DrawUniform(draw, region)
	}

	iterateLimit := renderer.app.BaseConfig().IterateLimit
	// Add detail from the small regions next
	for _, region := range smallRegions {
		RenderSequenceRegion(region, draw, iterateLimit)
	}

	return draw.Picture(), nil
}

func (renderer RegionRenderStrategy) SubdivideRegions(whole RegionNumerics) ([]RegionNumerics, []RegionNumerics) {
	// Lots of preallocated space for regions and region pointers
	completeRegions := make([]RegionNumerics, 0, base.AllocMedium)
	smallRegions := make([]RegionNumerics, 0, base.AllocMedium)
	splittingRegions := make([]RegionNumerics, 1, base.AllocMedium)
	regionConfig := renderer.app.RegionConfig()
	baseConfig := renderer.app.BaseConfig()
	collapseBound := int(regionConfig.CollapseSize)

	// Split regions
	splittingRegions[0] = whole
	for i := 0; i < len(splittingRegions); i++ {
		splitee := splittingRegions[i]

		splitee.ClaimExtrinsics()
		// There are three things that can happen to a region...
		//
		// A. The region can be so small that we divide no further
		if Collapse(splitee, collapseBound) {

			smallRegions = append(smallRegions, splitee)
		} else {
			// If the region is not too small, two things can happen
			// B. The region needs subdivided because it covers distinct parts of the plane
			if Subdivide(splitee, baseConfig.IterateLimit, regionConfig.GlitchSamples) {
				splittingRegions = append(splittingRegions, splitee.Children()...)
				// C. The region need not be divided
			} else {
				completeRegions = append(completeRegions, splitee)
			}
		}
	}

	return completeRegions, smallRegions
}
