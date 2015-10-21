package libgodelbrot

type BigSequentialNumerics struct {
    BigBaseNumerics
    sequencer func(i int, j int, member BigMandelbrotMember)
    members []PixelMember
}

func CreateBigSequentialNumerics(base BigBaseNumerics) BigSequentialNumerics{
    numerics := BigSequentialNumerics{
        BigBaseNumerics: base,
        members: make([]PixelMember, 0, allocTiny),
    }
}

func (bigFloat *BigSequentialNumerics) MandelbrotSequence(iterLimit uint8) {
    topLeft := bigFloat.PlaneTopLeft()

    imageLeft, imageTop := big.PictureMin()
    imageRight, imageBottom := big.PictureMax()
    rUnit, iUnit := big.PixelSize()
    divergeLimit := big.DivergeLimit()

    x := topLeft.Real()
    for i := imageLeft; i < imageRight; i++ {
        y := topLeft.Imag()
        for j := imageTop; j < imageBottom; j++ {
            member := BigMandelbrotMember{
                C: BigComplex{x, y},
            }
            member.Mandelbrot(iterLimit)
            bigFloat.Sequencer(i, j, member)
            y.Sub(y, iUnit)
        }
        x.Add(x, rUnit)
    }
}

func (big *BigSequentialNumerics) ImageDrawSequencer(draw DrawingContext) {
    big.sequencer = func (i, j int, member BigMandelbrotMember) {
        draw.DrawPointAt(i, j, member)
    }
}

func (big *BigSequentialNumerics) MemberCaptureSequencer() {
    big.sequencer = func (i, j int, member BigMandelbrotMember) {
        big.members = append(big.members, PixelMember{
            I: i, 
            J: j,
            MandelbrotMember: member.MandelbrotMember,
        }
    }
}

func (big *BigSequentialNumerics) CapturedMembers() []PixelMember {
    return big.members
}