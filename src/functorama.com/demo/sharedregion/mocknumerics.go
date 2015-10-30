package sharedregion

import (
    "functorama.com/demo/region"
    "functorama.com/demo/sequence"
)

type MockThreadPrototype struct {
    TGrabThreadPrototype bool
}

func (mock *MockThreadPrototype) GrabThreadPrototype(threadId uint) {
    mock.TGrabThreadPrototype = true
}

type MockSequence struct {
    sequence.MockNumerics
    MockThreadPrototype
}

type MockNumerics struct {
    region.MockNumerics
    MockThreadPrototype
    TSharedChildren  bool
    TSharedRegionSequence bool

    ShareNext []SharedRegionNumerics
    SharedSequence SharedSequenceNumerics
}

func (mock *MockNumerics) SharedChildren() []SharedRegionNumerics {
    mock.TSharedChildren = true
    return mock.ShareNext
}

func (mock *MockNumerics) SharedRegionSequence() SharedSequenceNumerics {
    mock.TSharedRegionSequence = true
    return mock.SharedSequence
}
