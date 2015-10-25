package libgodelbrot

type SequentialNumericsFactory interface {
    Sequence() SequentialNumerics
}

// SequentialNumerics provides sequential (column-wise) rendering calculations
type SequentialNumerics interface {
    OpaqueFlyweightProxy
    MandelbrotSequence(iterateLimit uint8)
    ImageDrawSequencer(draw DrawingContext)
    MemberCaptureSequencer()
    CapturedMembers() []PixelMember
    SubImage(rect image.Rectangle)
}