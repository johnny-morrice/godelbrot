package nativebase

import (
    "github.com/johnny-morrice/godelbrot/internal/base"
)

type NativeCoordProvider interface {
    NativeUserCoords() (complex128, complex128)
}

type RenderApplication interface {
    base.RenderApplication
    NativeCoordProvider
}