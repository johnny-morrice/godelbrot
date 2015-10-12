package libgodelbrot

type NativeSequentialNumerics struct {
    NativeBaseNumerics
    sequencer func(i int, j int, member NativeMandelbrotMember)
    members []PixelMember
}

func (native *NativeSequentialNumerics) MandelbrotSequence() {
    topLeft := native.PlaneTopLeft()

    imageLeft, imageTop := Native.PictureMin()
    imageRight, imageBottom := Native.PictureMax()
    rUnit, iUnit := Native.PlaneUnits()

    x := real(topLeft)
    for i := imageLeft; i < imageRight; i++ {
        y := imag(topLeft)
        for j := imageTop; j < imageBottom; j++ {
            member := NativeMandelbrotMember{
                C: complex(x, y),
            }
            &member.Mandelbrot(native.IterateLimit(), native.DivergeLimit())
            native.sequencer(i, j, *member)
            y -= iUnit
        }
        x += rUnit
    }
}

func (native *NativeSequentialNumerics) ImageDrawSequencer() {
    native.sequencer = DrawPointAt
}

func (native *NativeSequentialNumerics) MemberCaptureSequencer() {
    native.Sequencer = func (i, j int, member NativeMandelbrotMember) {
        native.members = append(native.members, PixelMember{
            I: i, 
            J: j,
            MandelbrotMember: NativeMandelbrotMember.MandelbrotMember,
        }
    }
}

func (native *NativeSequentialNumerics) CapturedMembers() []PixelMember {
    return native.members
}