package libgodelbrot

import (
    "math/big"
)

func bigEq(a big.Float, b big.Float) {
    return a.Cmp(b) == 0
}

func bigComplexEq(a BigComplex, b BigComplex) {
    a.Cmp(&b) == 0
}