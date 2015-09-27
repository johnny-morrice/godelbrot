package libgodelbrot

import (
    "image"
    "math"
)

// User input 
type RenderParameters struct {
    IterateLimit uint8
    DivergeLimit float64
    Width uint
    Height uint
    XOffset float64
    YOffset float64
    Zoom float64
    RegionCollapse uint
}

// Machine prepared input, caching interim results
type RenderConfig struct {
    RenderConfig
    // One pixel's space on the plane
    HorizUnit float64
    VerticalUnit float64
    ImageLeft uint
    ImageTop uint
}

func (args RenderParameters) Configure() &RenderParameters {
    size := args.PlaneTopLeft() - args.PlaneBottomRight()
    config := RenderParameters{args}
    config.HorizUnit := real(size) / float64(config.Width)
    config.VerticalUnit := imag(size) / float64(config.Height)
    return &config
}

// Top left of window onto complex plane
func (config RenderParameters) PlaneTopLeft() complex128 {
    return complex(config.XOffset, config.YOffset)
}

// Top right of window onto complex plane
func (config RenderParameters) PlaneBottomRight() complex128 {
    return windowSize * complex(config.Zoom, 0)
}

func (config RenderConfig) ImageTopLeft() (uint, uint) {
    return config.ImageLeft, config.ImageTop
}

func (config RenderParameters) BlankImage() image.NRGBA {
    pic := image.NewNRGBA(image.Rectangle{
        Min: image.ZP,
        Max: image.Point{
            X: widthI,
            Y: heightI,
        },
    })
}

func (config RenderConfig) PlaneToPixel(c complex128) (uint, uint) {
    // translate before scale
    x := (real(c) - config.XOffset) / config.HorizUnit
    y := (imag(c) - config.YOffset) / config.VerticalUnit
    // Remember that we draw downwards
    return math.Floor(x), math.Ceil(-y)
}

func (config RenderConfig) RegionRect(region *Region) image.Rectangle {
    l, t := config.PlaneToPixel(region.topLeft.c)
    r, b := config.PlaneToPixel(region.bottomRight.c)
    return image.Rect(int(l), int(t), int(r), int(b))
}