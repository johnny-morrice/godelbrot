package libgodelbrot

import (
    "image/draw"
    "image"
)

func RegionRender(configP *RenderConfig, palette Palette) (*image.NRGBA, error) {
    config := *configP
    initialRegion := NewRegion(config.WindowTopLeft(), config.WindowBottomRight())

    uniformRegions, smallRegions := subdivideRegions(initialRegion)

    // Render image
    pic := config.BlankImage()
    // Draw uniform regions first
    for _, region := range uniformRegions {
        member := region.midPoint.membership
        uniform := image.NewUniform(palette.Color(member))
        draw.Draw(pic, config.RegionRect(region), uniform, image.ZP, draw.Src)
    }

    // Add detail from the small regions next
    for _, region := range smallRegions {
        // Create config for rendering this region
        smallConfig := RenderConfig{config}
        regionConfig(region, configP, &smallConfig)
        SequentialRenderImage(smallConfig, palette, pic)
    }
    return pic, nil
}

func subdivideRegions(config *RenderConfig, whole *Region) ([]*Region, []*Region) {
   // Lots of preallocated space for regions and region pointers
    const meg uint = 1048576
    completeRegions := make([]*Region, 0, meg)
    smallRegions := make([]*Region, 0, meg)
    splittingRegions := make([]*Region, 1, meg)
    iterateLimit := config.IterateLimit
    divergeLimit := config.DivergeLimit

    // Split regions
    splittingRegions[0] = initialRegion
    for i := 0; i < len(splittingRegions); i++ {
        splitee := splittingRegions[i]
        x, y := splitee.PixelSize()
        // There are three things that can happen to a region...
        //
        // A. The region can be so small that we divide no further
        if x <= config.RegionCollapse || y <= config.RegionCollapse {
            smallRegions = append(smallRegions, splitee)
        } else {
            // If the region is not too small, two things can happen
            subregion := splitee.Subdivide(iterateLimit, divergeLimit)
            // B. The region needs subdivided because it covers chaoticall distinct parts of the plane 
            if subregion.populated {
                splittingRegions = append(splittingRegions, subregion.children...)
                // C. The region 
            } else {
                completeRegions = append(completeRegions, splitee)
            }
        }
    }

    return completeRegions, smallRegions
}

// Write image and plane position data to the small config
func regionConfig(smallRegion *Region, largeConfig *RenderConfig, smallConfig *RenderConfig) {
    rect := largeConfig.RegionRect(smallRegion)
    topLeft = smallRegion.topLeft.c
    smallConfig.Width = uint(rect.Dx)
    smallConfig.Height = uint(rect.Dy)
    smallConfig.ImageLeft = uint(rect.Min.X)
    smallConfig.ImageTop = uint(rect.Max.Y)
    smallConfig.XOffset = real(topLeft)
    smallConfig.YOffset = imag(topLeft)
}