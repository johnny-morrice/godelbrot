package bigbase

import (
    "functorama.com/demo/base"
)

type BigCoordProvider interface {
    BigUserCoords() (*BigComplex, *BigComplex)
    Precision() uint
}

type RenderApplication interface {
    base.RenderApplication
    BigCoordProvider
}