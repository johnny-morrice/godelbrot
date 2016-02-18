package sequence

import (
    "github.com/johnny-morrice/godelbrot/base"
    "github.com/johnny-morrice/godelbrot/draw"
)

type MockRenderApplication struct {
    base.MockRenderApplication
    draw.MockContextProvider

    TSequenceNumericsFactory bool

    SequenceFactory SequenceNumericsFactory
}

func (mock *MockRenderApplication) SequenceNumericsFactory() SequenceNumericsFactory {
    mock.TSequenceNumericsFactory = true
    return mock.SequenceFactory
}

type MockFactory struct {
    TBuild bool
    Numerics SequenceNumerics
}

func (mock *MockFactory) Build() SequenceNumerics {
    mock.TBuild = true
    return mock.Numerics
}