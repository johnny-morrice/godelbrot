package libgodelbrot

import (
    "image"
)

func SequentialRender(config *RenderConfig, palette Palette) (*image.NRGBA, error) {
    pic := config.BlankImage()
    SequentialRenderImage(config, palette, pic)
    return pic, nil
}

func SequentialRenderImage(config *RenderConfig, palette Palette, pic *image.NRGBA) {
    topLeft := config.PlaneTopLeft()

    widthI := int(config.Width)
    heightI := int(config.Height)

    imageLeft, imageTop := config.ImageTopLeft()

    x := real(topLeft)
    y := imag(topLeft)
    for i := int(imageLeft); i < widthI; i++ {
        y = imag(topLeft)
        for j := int(imageTop); j < heightI; j++ {
            c := complex(x, y)
            member := Mandelbrot(c, config.IterateLimit, config.DivergeLimit)
            color := palette.Color(member)
            pic.Set(i, j, color)
            y -= config.VerticalUnit
        }
        x += config.HorizUnit
    }
}