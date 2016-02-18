package region

import (
	"image"
	"github.com/johnny-morrice/godelbrot/base"
	"github.com/johnny-morrice/godelbrot/draw"
)

type RegionRenderStrategy struct {
	factory RegionNumericsFactory
	context draw.DrawingContext
	regionConfig RegionConfig
}

func Make(app RenderApplication) *RegionRenderStrategy {
	return &RegionRenderStrategy{
		factory: app.RegionNumericsFactory(),
		context: app.DrawingContext(),
		regionConfig: app.RegionConfig(),
	}
}

// The RegionRenderStrategy implements RenderNumerics with this method that
// draws the Mandelbrot set uses a "similar rectangles" optimization
func (renderer RegionRenderStrategy) Render() (*image.NRGBA, error) {
	// The numerics system is by default a region covering the whole image
	initialRegion := renderer.factory.Build()
	uniformRegions, smallRegions := renderer.SubdivideRegions(initialRegion)

	// Draw uniform regions first
	for _, region := range uniformRegions {
		region.ClaimExtrinsics()
		DrawUniform(renderer.context, region)
	}

	// Add detail from the small regions next
	for _, region := range smallRegions {
		RenderSequenceRegion(region, renderer.context)
	}

	return renderer.context.Picture(), nil
}

func (renderer RegionRenderStrategy) SubdivideRegions(whole RegionNumerics) ([]RegionNumerics, []RegionNumerics) {
	// Lots of preallocated space for regions and region pointers
	completeRegions := make([]RegionNumerics, 0, base.AllocMedium)
	smallRegions := make([]RegionNumerics, 0, base.AllocMedium)
	splittingRegions := make([]RegionNumerics, 1, base.AllocMedium)
	collapseBound := int(renderer.regionConfig.CollapseSize)

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
			divided := Subdivide(splitee)
			if divided {
				splittingRegions = append(splittingRegions, splitee.Children()...)
				// C. The region need not be divided
			} else {
				completeRegions = append(completeRegions, splitee)
			}
		}
	}

	return completeRegions, smallRegions
}
