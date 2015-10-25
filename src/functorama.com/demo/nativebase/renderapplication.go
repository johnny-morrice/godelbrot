package nativebase

import (
    "functorama.com/demo/base"
)

type RenderApplication interface {
    base.RenderApplication
    NativeUserCoords() (complex128, complex128)
}