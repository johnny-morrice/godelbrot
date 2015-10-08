package libgodelbrot

func (config *NativeConfig) MandelbrotSequence() {
    topLeft := config.PlaneTopLeft()

    imageLeft, imageTop := config.ImageTopLeft()
    maxH := int(config.Width) + int(imageLeft)
    maxV := int(config.Height) + int(imageTop)

    x := real(topLeft)
    for i := int(imageLeft); i < maxH; i++ {
        y := imag(topLeft)
        for j := int(imageTop); j < maxV; j++ {
            member := NativeMandelbrotMember{
                C: complex(x, y),
            }
            &member.Mandelbrot(config.IterateLimit, config.DivergeLimit)
            config.Sequencer(i, j, *member)
            y -= config.VerticalUnit
        }
        x += config.HorizUnit
    }
}