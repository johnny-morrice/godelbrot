package bigbase

import (
	"fmt"
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

func (c *BigComplex) Real() *big.Float {
	return &c.R
}

func (c *BigComplex) Imag() *big.Float {
	return &c.I
}

func (c *BigComplex) Add(q *BigComplex, u *BigComplex) {
	c.Real().Add(q.Real(), u.Real())
	c.Imag().Add(q.Imag(), u.Imag())
}

// Create a new Float, supplying a precision
func CreateBigFloat(f float64, prec uint) big.Float {
	b := big.Float{}
	b.SetFloat64(f)
	b.SetPrec(prec)
	return b
}

// DbgF is a big.Float with an easy to grok string representation
type DbgF big.Float

func (df DbgF) String() string {
	bf := big.Float(df)
	message := "Val: %v Prec: %v"
	val, _ := bf.Float64()
	return fmt.Sprintf(message, val, bf.Prec())
}

// DbgC is a BigComplex with an easy to grok string representation
type DbgC BigComplex

func (dc DbgC) String() string {
	return fmt.Sprintf("BigComplex{%v, %v}", DbgF(dc.R), DbgF(dc.I))
}
