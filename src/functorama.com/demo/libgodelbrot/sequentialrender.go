package libgodelbrot

// Normalized size of window onto complex plane
const windowSize complex128 = 1 + 1i

type RenderParameters struct {
    IterateLimit uint
    DivergeLimit float64
    Width uint
    Height uint
    XOffset float64
    YOffset float64
    Zoom float64
}

type SequentialRenderer interface {}

func NewSequentialRenderer() SequentialRenderer {
    return RenderParameters{}
}

func (renderer *SequentialRenderer) Render(args RenderParameters) image.NRGBA {
    var bottomRight complex128 = windowSize * zoom
    horizUnit := real(bottomRight) / RenderParameters.Width
    verticalUnit := imag(bottomRight) / RenderParameters.Height

    pic := image.NewNRGBA(image.Rectangle{
        Min: image.ZP,
        Max: image.Point{
            X: RenderParameters.Width,
            Y: RenderParameters.Height,
        },
    })

    x := 0
    for i := 0; i < RenderParameters.Width; i++; {
        y := 0
        for j := 0; j < RenderParameters.Height; j++ {
            c := complex(x, y)
            member := Mandelbrot(c, RenderParameters.IterateLimit, RenderParameters.DivergeLimit)
            color := MandelbrotColor(member)
            pic.Set(i, j, color)
            y += verticalUnit
        }
        x += horizUnit
    }
    
    return image
}