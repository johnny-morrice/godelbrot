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

    imageLeft, imageTop := config.ImageTopLeft()
    maxH := int(config.Width) + int(imageLeft)
    maxV := int(config.Height) + int(imageTop)

    if imageLeft == 0 && imageTop == 0 {
        _ = "breakpoint"
    }

    x := real(topLeft)
    for i := int(imageLeft); i < maxH; i++ {
        y := imag(topLeft)
        for j := int(imageTop); j < maxV; j++ {
            c := complex(x, y)
            member := Mandelbrot(c, config.IterateLimit, config.DivergeLimit)
            color := palette.Color(member)
            pic.Set(i, j, color)
            y -= config.VerticalUnit
        }
        x += config.HorizUnit
    }
}