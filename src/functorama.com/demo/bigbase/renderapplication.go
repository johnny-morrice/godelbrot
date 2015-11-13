package bigbase

import (
    "functorama.com/demo/base"
)

type RenderApplication interface {
    base.RenderApplication
    BigUserCoords() (*BigComplex, *BigComplex)
    Precision() uint
}