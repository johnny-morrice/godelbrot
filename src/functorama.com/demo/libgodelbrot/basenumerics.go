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

// Numerics that are aware of a picture and of the Mandelbrot iteration limit
type BaseNumerics struct {
    picXMin int
    picXMax int // exclusive maximum
    picYMin int
    picYMax int // exclusive maximum
}

func CreateBaseNumerics(app RenderApplication) BaseNumerics {
    iLimit, _ := app.Limits()
    return BaseNumerics{
        picXMin: 0,
        picXMax: app.PictureWidth(),
        picYMin: 0,
        picYMax: app.PictureHeight(),
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