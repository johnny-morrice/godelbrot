package bigbase

import (
    "functorama.com/demo/base"
)

type MockRenderApplication struct {
    base.MockRenderApplication

    TBigUserCoords bool
    TPrecision bool

    UserMin BigComplex
    UserMax BigComplex
    Prec uint
}

var _ RenderApplication = (*MockRenderApplication)(nil)

func (mra *MockRenderApplication) Precision() uint {
    mra.TPrecision = true
    return mra.Prec
}

func (mra *MockRenderApplication) BigUserCoords() (*BigComplex,*BigComplex) {
    mra.TBigUserCoords = true
    return &mra.UserMin, &mra.UserMax
}