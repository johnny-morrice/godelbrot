package libgodelbrot

type NativeSequentialNumerics struct {
	NativeBaseNumerics
	sequencer func(i int, j int, member NativeMandelbrotMember)
	members   []PixelMember
}

func (native *NativeSequentialNumerics) MandelbrotSequence(iterLimit uint8) {
	topLeft := native.PlaneTopLeft()

	imageLeft, imageTop := native.PictureMin()
	imageRight, imageBottom := native.PictureMax()
	rUnit, iUnit := native.PixelSize()

	x := real(topLeft)
	for i := imageLeft; i < imageRight; i++ {
		y := imag(topLeft)
		for j := imageTop; j < imageBottom; j++ {
			member := NativeMandelbrotMember{
				C: complex(x, y),
			}
			member.Mandelbrot(iterLimit)
			native.sequencer(i, j, *member)
			y -= iUnit
		}
		x += rUnit
	}
}

func (native *NativeSequentialNumerics) ImageDrawSequencer(draw DrawingContext) {
	native.sequencer = func(i, j int, member NativeMandelbrotMember) {
		draw.DrawPointAt(PixelMember{i, j, member.MandelbrotMember})
	}
}

func (native *NativeSequentialNumerics) MemberCaptureSequencer() {
	native.Sequencer = func(i, j int, member NativeMandelbrotMember) {
		native.members = append(native.members, PixelMember{
			I:      i,
			J:      j,
			Member: member,
		})
	}
}

func (native *NativeSequentialNumerics) CapturedMembers() []PixelMember {
	return native.members
}
