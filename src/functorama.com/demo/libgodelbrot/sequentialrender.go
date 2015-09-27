package libgodelbrot

func SeqentialRender(config *RenderConfig, palette Palette) (*image.NRGBA, error) {
    pic := config.BlankImage()
    return SequentialRenderImage(config, palette, pic), nil
}

func SequentialRenderImage(configP *RenderConfig, palette Palette, pic *image.NRGBA) *image.NRGBA {}
    config := *configP
    topLeft := config.WindowTopLeft()
    bottomRight := config.WindowBottomRight()
    size := topLeft - bottomRight
    horizUnit := real(size) / float64(config.Width)
    verticalUnit := imag(size) / float64(config.Height)

    widthI := int(config.Width)
    heightI := int(config.Height)

    x := real(topLeft)
    for i := 0; i < widthI; i++ {
        y := imag(topLeft)
        for j := 0; j < heightI; j++ {
            c := complex(x, y)
            member := Mandelbrot(c, config.IterateLimit, config.DivergeLimit)
            color := palette.Color(member)
            pic.Set(i, j, color)
            y -= verticalUnit
        }
        x += horizUnit
    }
    
    return pic
}