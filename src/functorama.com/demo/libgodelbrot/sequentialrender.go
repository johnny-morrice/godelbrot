package libgodelbrot

import (
    "image"
)

func SequentialRender(config *RenderConfig, palette Palette) (*image.NRGBA, error) {
    pic := config.BlankImage()
    SequentialRenderImage(CreateContext(config, palette, pic))
    return pic, nil
}

func SequentialRenderImage(drawingContext DrawingContext) {
    MandelbrotSequence(drawingContext.Config, drawingContext.DrawPointAt)
}

type Sequencer func(i int, j int, member MandelbrotMember)

func MandelbrotSequence(config *RenderConfig, sequencer Sequencer) {
    topLeft := config.PlaneTopLeft()

    imageLeft, imageTop := config.ImageTopLeft()
    maxH := int(config.Width) + int(imageLeft)
    maxV := int(config.Height) + int(imageTop)

    x := real(topLeft)
    for i := int(imageLeft); i < maxH; i++ {
        y := imag(topLeft)
        for j := int(imageTop); j < maxV; j++ {
            c := complex(x, y)
            member := Mandelbrot(c, config.IterateLimit, config.DivergeLimit)
            sequencer(i, j, member)
            y -= config.VerticalUnit
        }
        x += config.HorizUnit
    }
}