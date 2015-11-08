package sharedregion

import (
    "functorama.com/demo/region"
)

type MockThreadPrototype struct {
    TGrabThreadPrototype bool
}

func (mock *MockThreadPrototype) GrabWorkerPrototype(threadId uint16) {
    mock.TGrabThreadPrototype = true
}

type MockSequence struct {
    region.MockProxySequence
    MockThreadPrototype
}

type MockNumerics struct {
    region.MockNumerics
    MockThreadPrototype
    TSharedChildren  bool
    TSharedRegionSequence bool

    SharedMockChildren []*MockNumerics
    SharedMockSequence *MockSequence
}

func (mock *MockNumerics) SharedChildren() []SharedRegionNumerics {
    mock.TSharedChildren = true
    abstract := make([]SharedRegionNumerics, len(mock.SharedMockChildren))

    for i, child := range mock.SharedMockChildren {
        abstract[i] = child
    }

    return abstract
}

func (mock *MockNumerics) SharedRegionSequence() SharedSequenceNumerics {
    mock.TSharedRegionSequence = true
    return mock.SharedMockSequence
}