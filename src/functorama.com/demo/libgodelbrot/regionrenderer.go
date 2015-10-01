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
    initialRegion := WholeRegion(config)
    escapePointHeap := NewEscapePointHeap(K64)
    renderConfigHeap := NewRenderConfigHeap(config, K64)
    uniformRegions, smallRegions := subdivideRegions(config, initialRegion, escapePointHeap)

    // Draw uniform regions first
    for _, region := range uniformRegions {
        drawingContext.DrawUniform(region)
    }

    // Add detail from the small regions next
    for _, region := range smallRegions {
        RenderSequentialRegion(region, drawingContext, renderConfigHeap)
    }
}

func subdivideRegions(config *RenderConfig, whole Region, heap *EscapePointHeap) ([]Region, []Region) {
   // Lots of preallocated space for regions and region pointers
    completeRegions := make([]Region, 0, K64)
    smallRegions := make([]Region, 0, K64)
    splittingRegions := make([]Region, 1, K64)

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
            subregion := splitee.Subdivide(config, heap)
            // B. The region needs subdivided because it covers distinct parts of the plane 
            if subregion.populated {
                splittingRegions = append(splittingRegions, subregion.children...)
                // C. The region need not be divided
            } else {
                completeRegions = append(completeRegions, splitee)
            }
        }
    }

    return completeRegions, smallRegions
}

func RenderSequentialRegion(region Region, drawingContext DrawingContext, heap *RenderConfigHeap) {
    // Create config for rendering this region
    smallConfig := heap.Subconfig(region)
    smallContext := CreateContext(smallConfig, drawingContext.ColorPalette, drawingContext.Pic)
    SequentialRenderImage(smallContext)
}