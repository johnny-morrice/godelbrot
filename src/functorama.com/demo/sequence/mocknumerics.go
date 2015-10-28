package sequence

import (
    "image"
    "functorama.com/demo/base"
    "functorama.com/demo/draw"
)

type MockNumerics struct {
    TMandelbrotSequence bool
    TImageDrawSequencer bool
    TMemberCaptureSequencer bool
    TCapturedMembers bool
    TSubImage bool

    Captured []base.PixelMember
}

func (mock *MockNumerics) MandelbrotSequence(iterateLimit uint8) {
    mock.TMandelbrotSequence = true
}

func (mock *MockNumerics) ImageDrawSequencer(context draw.DrawingContext) {
    mock.TImageDrawSequencer = true
}

func (mock *MockNumerics) MemberCaptureSequencer() {
    mock.TMemberCaptureSequencer = true
}

func (mock *MockNumerics) CapturedMembers() []base.PixelMember {
    mock.TCapturedMembers = true
    return mock.Captured
}

func (mock *MockNumerics) SubImage(rect image.Rectangle) {
    mock.TSubImage = true
}
