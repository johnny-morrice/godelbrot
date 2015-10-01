package libgodelbrot

import (
    "image"
    "math"
    "fmt"
)

// Co-ordinate frames
const (
    CornerFrame = iota
    ZoomFrame = iota
)

// User input 
type RenderParameters struct {
    IterateLimit uint8
    DivergeLimit float64
    Width uint
    Height uint
    Zoom float64
    RegionCollapse uint
    // Co-ordinate frames
    Frame uint
    // Top left of view onto plane
    TopLeft complex128
    // Optional Bottom right corner
    BottomRight complex128
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
        TopLeft: MagicOffset,
        Zoom: DefaultZoom,
        Frame: ZoomFrame,
        RegionCollapse: DefaultCollapse,
    }
    return params.Configure()
}

func (args RenderParameters) PlaneSize() complex128 {
    if args.Frame == ZoomFrame {
        return complex(args.Zoom, 0) * MagicSetSize
    } else if args.Frame == CornerFrame {
        tl := args.TopLeft
        br := args.BottomRight
        return complex(real(br) - real(tl), imag(tl) - imag(br))
    } else {
        args.framePanic()
    }
    panic("Bug")
    return 0
}

func (args RenderParameters) Configure() *RenderConfig {
    planeSize := args.PlaneSize()
    planeWidth := real(planeSize)
    planeHeight := imag(planeSize)
    return &RenderConfig{
        RenderParameters: args,
        HorizUnit: planeWidth / float64(args.Width),
        VerticalUnit: planeHeight / float64(args.Height),
        ImageLeft: 0,
        ImageTop: 0,
    }
}

func (config RenderParameters) PlaneTopLeft() complex128 {
    return config.TopLeft
}

// Top right of window onto complex plane
func (config RenderParameters) PlaneBottomRight() complex128 {
    if config.Frame == ZoomFrame {
        scaled := MagicSetSize * complex(config.Zoom, 0)
        topLeft := config.PlaneTopLeft()
        right := real(topLeft) + real(scaled)
        bottom := imag(topLeft) - imag(scaled)
        return complex(right, bottom)
    } else if config.Frame == CornerFrame {
        return config.BottomRight
    } else {
        config.framePanic()
    }
    panic("Bug")
    return 0
}

func (config RenderParameters) framePanic() {
    panic(fmt.Sprintf("Unknown frame: %v", config.Frame))
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

func (config RenderConfig) PlaneToPixel(c complex128) (rx uint, ry uint) {
     // Translate x
    tx := real(c) - real(config.TopLeft)
    // Scale x
    sx := tx / config.HorizUnit

    // Translate y
    ty := imag(c) - imag(config.TopLeft) 
    // Scale y
    sy := ty / config.VerticalUnit
    
    rx = uint(math.Floor(sx))
    // Remember that we draw downwards
    ry = uint(math.Ceil(-sy))

    return
}