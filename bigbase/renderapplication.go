package bigbase

import (
    "github.com/johnny-morrice/godelbrot/base"
)

type BigCoordProvider interface {
    BigUserCoords() (*BigComplex, *BigComplex)
    Precision() uint
}

type RenderApplication interface {
    base.RenderApplication
    BigCoordProvider
}