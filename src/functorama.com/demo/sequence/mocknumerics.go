package sequence

import (
    "image"
    "functorama.com/demo/base"
)

type MockNumerics struct {
    TSequence bool
    TArea bool
    TSubImage bool
}

// Check MockNumerics implements SequenceNumerics interface
var _ SequenceNumerics = (*MockNumerics)(nil)

func (mn *MockNumerics) Sequence(iterateLimit uint8) <-chan base.PixelMember {
    mn.TSequence = true

    out := make(chan base.PixelMember, 1)
    close(out)
    return out
}

func (mn *MockNumerics) SubImage(rect image.Rectangle) {
    mn.TSubImage = true
}

func (mn *MockNumerics) Area() int {
    mn.TArea = true
    return 1
}
