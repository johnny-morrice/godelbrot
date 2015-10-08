package libgodelbrot

func (config *BigConfig) MandelbrotSequence() {
    topLeft := config.PlaneTopLeft()

    imageLeft, imageTop := config.ImageTopLeft()
    maxH := int(config.Width) + int(imageLeft)
    maxV := int(config.Height) + int(imageTop)

    x := topLeft.Real()
    for i := int(imageLeft); i < maxH; i++ {
        y := topLeft.Imag()
        for j := int(imageTop); j < maxV; j++ {
            member := BigMandelbrotMember{
                C: BigComplex{x, y},
            }
            &member.Mandelbrot(config.IterateLimit, config.DivergeLimit)
            config.Sequencer(i, j, member)
            y.Sub(y, config.VerticalUnit)
        }
        x.Add(x, config.HorizUnit)
    }
}