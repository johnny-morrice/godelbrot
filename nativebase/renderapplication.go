package nativebase

import (
    "functorama.com/demo/base"
)

type NativeCoordProvider interface {
    NativeUserCoords() (complex128, complex128)
}

type RenderApplication interface {
    base.RenderApplication
    NativeCoordProvider
}