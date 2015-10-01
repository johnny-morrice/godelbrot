package libgodelbrot

import (
    "image/draw"
    "image"
)

func RegionRender(config *RenderConfig, palette Palette) (*image.NRGBA, error) {
    pic := config.BlankImage()
    RegionRenderImage(config, palette, pic)
    return pic, nil
}

func RegionRenderImage(config *RenderConfig, palette Palette, pic *image.NRGBA) {
    initialRegion := NewRegion(config.PlaneTopLeft(), config.PlaneBottomRight())
    uniformRegions, smallRegions := subdivideRegions(config, initialRegion)

    // Draw uniform regions first
    for _, region := range uniformRegions {
        member := region.midPoint.membership
        color := palette.Color(member)
        uniform := image.NewUniform(color)
        rect := region.Rect(config)
        draw.Draw(pic, rect, uniform, image.ZP, draw.Src)
    }

    // Add detail from the small regions next
    for _, region := range smallRegions {
        // Create config for rendering this region
        smallConfig := *config
        regionConfig(region, config, &smallConfig)
        SequentialRenderImage(&smallConfig, palette, pic)
    }
}

func subdivideRegions(config *RenderConfig, whole *Region) ([]*Region, []*Region) {
   // Lots of preallocated space for regions and region pointers
    const meg uint = 1048576
    completeRegions := make([]*Region, 0, meg)
    smallRegions := make([]*Region, 0, meg)
    splittingRegions := make([]*Region, 1, meg)

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
            subregion := splitee.Subdivide(config)
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

var badPtr *Region

// Write image and plane position data to the small config
func regionConfig(smallRegion *Region, largeConfig *RenderConfig, smallConfig *RenderConfig) {
    rect := smallRegion.Rect(largeConfig)
    smallConfig.Width = uint(rect.Dx())
    smallConfig.Height = uint(rect.Dy())
    smallConfig.ImageLeft = uint(rect.Min.X)
    smallConfig.ImageTop = uint(rect.Min.Y)
    smallConfig.TopLeft = smallRegion.topLeft.c
    smallConfig.BottomRight = smallRegion.bottomRight.c
    smallConfig.Frame = CornerFrame
}