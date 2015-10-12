package libgodelbrot

func (bigFloat *BigFloatNumerics) MandelbrotSequence() {
    topLeft := bigFloat.PlaneTopLeft()

    imageLeft, imageTop := bigFloat.ImageTopLeft()
    maxH := int(bigFloat.Width) + int(imageLeft)
    maxV := int(bigFloat.Height) + int(imageTop)

    x := topLeft.Real()
    for i := int(imageLeft); i < maxH; i++ {
        y := topLeft.Imag()
        for j := int(imageTop); j < maxV; j++ {
            member := BigMandelbrotMember{
                C: BigComplex{x, y},
            }
            &member.Mandelbrot(bigFloat.IterateLimit, bigFloat.DivergeLimit)
            bigFloat.Sequencer(i, j, member)
            y.Sub(y, bigFloat.VerticalUnit)
        }
        x.Add(x, bigFloat.HorizUnit)
    }
}