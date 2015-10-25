package bigbase

import (
	"math/big"
)

type BigComplex struct {
	R big.Float
	I big.Float
}

func CreateBigComplex(real float64, imag float64, prec uint) BigComplex {
	bigReal := CreateBigFloat(real, prec)
	bigImag := CreateBigFloat(imag, prec)
	return BigComplex{bigReal, bigImag}
}

func (c BigComplex) Real() FloatKind {
	return c.R
}

func (c BigComplex) Imag() FloatKind {
	return c.I
}

func (c BigComplex) Add(other BigComplex) {
	c.R.Add(other.R)
	c.I.Add(other.I)
}

// Use when you can assume accuracy is okay
func (f *big.Float) Float() float64 {
	native, acc := f.Float64()
	return native
}

// Create a new Float, supplying a precision
func CreateBigFloat(f float64, prec uint) *big.Float {
	b := big.Float{}
	b.Set(f)
	b.SetPrec(prec)
	return b
}
