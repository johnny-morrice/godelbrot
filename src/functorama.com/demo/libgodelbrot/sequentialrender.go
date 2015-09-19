package libgodelbrot

import (
    "image"
)

// Normalized size of window onto complex plane
const windowSize complex128 = 1 + 1i

type RenderParameters struct {
    IterateLimit uint8
    DivergeLimit float64
    Width uint
    Height uint
    XOffset float64
    YOffset float64
    Zoom float64
}

type SequentialRenderer struct {}

func NewSequentialRenderer() *SequentialRenderer {
    return &SequentialRenderer{}
}

func (renderer *SequentialRenderer) Render(argP *RenderParameters) (*image.NRGBA, error) {
    args := *argP
    var bottomRight complex128 = windowSize * complex(args.Zoom, 0)
    horizUnit := real(bottomRight) / float64(args.Width)
    verticalUnit := imag(bottomRight) / float64(args.Height)


    widthI := int(args.Width)
    heightI := int(args.Height)

    pic := image.NewNRGBA(image.Rectangle{
        Min: image.ZP,
        Max: image.Point{
            X: widthI,
            Y: heightI,
        },
    })

    palette := NewRedscalePalette()
    x := 0.0
    for i := 0; i < widthI; i++ {
        y := 0.0
        for j := 0; j < heightI; j++ {
            c := complex(x, y)
            member := Mandelbrot(c, args.IterateLimit, args.DivergeLimit)
            color := palette.Lookup(member)
            pic.Set(i, j, color)
            y += verticalUnit
        }
        x += horizUnit
    }
    
    return pic, nil
}