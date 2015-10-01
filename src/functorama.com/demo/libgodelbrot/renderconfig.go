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
    RenderParameters
    // One pixel's space on the plane
    HorizUnit float64
    VerticalUnit float64
    ImageLeft uint
    ImageTop uint
}

// Use magic values to create default config
func DefaultConfig() *RenderConfig {
    params := RenderParameters{
        IterateLimit: DefaultIterations,
        DivergeLimit: DefaultDivergeLimit,
        Width: DefaultImageWidth,
        Height: DefaultImageHeight,
        XOffset: real(MagicOffset),
        YOffset: imag(MagicOffset),
        Zoom: DefaultZoom,
        RegionCollapse: DefaultCollapse,
    }
    return params.Configure()
}

func (args RenderParameters) Configure() *RenderConfig {
    size := args.PlaneBottomRight() - args.PlaneTopLeft()
    return &RenderConfig{
        RenderParameters: args,
        HorizUnit: real(size) / float64(args.Width),
        VerticalUnit: imag(size) / float64(args.Height),
        ImageLeft: 0,
        ImageTop: 0,
    }
}

// Top left of window onto complex plane
func (config RenderParameters) PlaneTopLeft() complex128 {
    return complex(config.XOffset, config.YOffset)
}

// Top right of window onto complex plane
func (config RenderParameters) PlaneBottomRight() complex128 {
    return config.PlaneTopLeft() + (MagicSetSize * complex(config.Zoom, 0))
}

func (config RenderConfig) ImageTopLeft() (uint, uint) {
    return config.ImageLeft, config.ImageTop
}

func (config RenderParameters) BlankImage() *image.NRGBA {
    return image.NewNRGBA(image.Rectangle{
        Min: image.ZP,
        Max: image.Point{
            X: int(config.Width),
            Y: int(config.Height),
        },
    })
}

func (config RenderConfig) PlaneToPixel(c complex128) (uint, uint) {
    // translate before scale
    x := (real(c) - config.XOffset) / config.HorizUnit
    y := (imag(c) - config.YOffset) / config.VerticalUnit
    // Remember that we draw downwards
    return uint(math.Floor(x)), uint(math.Ceil(-y))
}

func (config RenderConfig) RegionRect(region *Region) image.Rectangle {
    l, t := config.PlaneToPixel(region.topLeft.c)
    r, b := config.PlaneToPixel(region.bottomRight.c)
    return image.Rect(int(l), int(t), int(r), int(b))
}