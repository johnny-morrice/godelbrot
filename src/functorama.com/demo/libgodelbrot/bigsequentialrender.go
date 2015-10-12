package libgodelbrot

func (bigFloat *BigFloatNumerics) MandelbrotSequence() {
    topLeft := bigFloat.PlaneTopLeft()

    imageLeft, imageTop := big.PictureMin()
    imageRight, imageBottom := big.PictureMax()
    rUnit, iUnit := big.PixelSize()
    iterLimit := native.IterateLimit()
    divergeLimit := native.DivergeLimit()

    x := topLeft.Real()
    for i := imageLeft; i < imageRight; i++ {
        y := topLeft.Imag()
        for j := imageTop; j < imageBottom; j++ {
            member := BigMandelbrotMember{
                C: BigComplex{x, y},
            }
            &member.Mandelbrot(iterLimit, divergeLimit)
            bigFloat.Sequencer(i, j, member)
            y.Sub(y, iUnit)
        }
        x.Add(x, rUnit)
    }
}