package libgodelbrot

import (
	"image"
)

func RegionRender(config *RenderConfig, palette Palette) (*image.NRGBA, error) {
	pic := config.BlankImage()
	RegionRenderImage(CreateContext(config, palette, pic))
	return pic, nil
}

func RegionRenderImage(drawingContext DrawingContext) {
	config := drawingContext.Config
	initialRegion := config.WholeRegionRenderContext()
	uniformRegions, smallRegions := SubdivideRegions(initialRegion)

	// Draw uniform regions first
	for _, region := range uniformRegions {
		drawingContext.DrawUniform(region)
	}

	// Add detail from the small regions next
	for _, region := range smallRegions {
		RenderSequentialRegion(region)
	}
}

func SubdivideRegions(whole RegionRenderContext) ([]RegionRenderContext, []RegionRenderContext) {
	// Lots of preallocated space for regions and region pointers
	completeRegions := make([]RegionRenderContext, 0, allocMedium)
	smallRegions := make([]RegionRenderContext, 0, allocMedium)
	splittingRegions := make([]RegionRenderContext, 1, allocMedium)

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

func RenderSequentialRegion(region RegionRenderContext) {
	smallContext := region.SubSequentialConfig()
	SequentialRenderImage(smallContext)
}
