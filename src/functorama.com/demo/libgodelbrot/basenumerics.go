package libgodelbrot

import (
    "image"
)

// A reusable notion of collapsable regions
type BaseRegionNumerics struct {
    glitchSamples int
    collapseSize int
}

func (collapse BaseRegionNumerics) CollapseSize() int {
    return collapse.collapseSize
}

// Number of samples to take when computing rendering glitch
func  (base BaseRegionNumerics) GlitchSamples() int {
    return base.glitchSamples
}

// Numerics that are aware of a picture and of the Mandelbrot iteration limit
type BaseNumerics struct {
    picXMin int
    picXMax int
    picYMin int
    picYMax int

    iterLimit int
}

func CreateBaseNumerics(render RenderApplication) BaseNumerics {
    iLimit, _ := render.Limits()
    return BaseNumerics{
        picXMin: 0,
        picXMax: pictureWidth,
        picYMin: 0,
        picYMax: pictureHeight,
        iterLimit: iLimit
    }
}

func (base *BaseNumerics) PictureMin() (int, int) {
    return base.picXMin, base.picYMin
}

func (base *BaseNumerics) PictureMax() (int, int) {
    return base.picXMax, base.picYMax
}


// Change the drawing context to a sub-part of the image
func (native *BaseNumerics) SubImage(rect image.Rectange) {
    native.picXMin = rect.Min.X
    native.picYmin = rect.Min.Y
    native.picXMax = rect.Max.X
    native.picYMax = rect.Max.Y    
}