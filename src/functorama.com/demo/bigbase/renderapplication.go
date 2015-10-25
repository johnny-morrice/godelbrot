package bigbase

import (
    "functorama.com/demo/base"
)

type BigRenderApplication struct {
    base.BaseRenderApplication
    BigUserCoords() (BigComplex, BigComplex)
}