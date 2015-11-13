package bigbase

import (
	"math/big"
)

func bigEq(a *big.Float, b *big.Float) bool {
	return a.Cmp(b) == 0
}

func bigComplexEq(a *BigComplex, b *BigComplex) bool {
	return bigEq(a.Real(), b.Real()) && bigEq(a.Imag(), b.Imag())
}
