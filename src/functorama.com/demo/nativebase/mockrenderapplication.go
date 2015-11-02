package nativebase

import (
    "functorama.com/demo/base"
)

type MockRenderApplication struct {
    base.MockRenderApplication

    TNativeUserCoords bool

    PlaneMin complex128
    PlaneMax complex128
    FixAspect bool
}

func (mock *MockRenderApplication) NativeUserCoords() (complex128, complex128) {
    mock.TNativeUserCoords = true
    return mock.PlaneMin, mock.PlaneMax
}