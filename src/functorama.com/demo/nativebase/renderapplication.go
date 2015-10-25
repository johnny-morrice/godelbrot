package nativebase

import (
    "functorama.com/demo/base"
)

type NativeRenderApplication interface {
    base.BaseRenderApplication
    NativeUserCoords() (complex128, complex128)
}