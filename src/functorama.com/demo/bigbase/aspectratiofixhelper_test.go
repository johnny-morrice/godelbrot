package bigbase

import (
	"math/big"
)

type aspectRatioFixHelper struct {
	pictureW uint
	pictureH uint

	rMin float64
	rMax float64
	iMin float64
	iMax float64

	expectRMin float64
	expectRMax float64
	expectIMin float64
	expectIMax float64
}

func (helper aspectRatioFixHelper) planeCoords() (BigComplex, BigComplex) {
	min := BigComplex{
		R: *big.NewFloat(helper.rMin),
		I: *big.NewFloat(helper.iMin),
	}
	max := BigComplex{
		R: *big.NewFloat(helper.rMax),
		I: *big.NewFloat(helper.iMax),
	}
	return min, max
}

func (helper aspectRatioFixHelper) expectCoords() (BigComplex, BigComplex) {
	min := BigComplex{
		R: *big.NewFloat(helper.expectRMin),
		I: *big.NewFloat(helper.expectIMin),
	}
	max := BigComplex{
		R: *big.NewFloat(helper.expectRMax),
		I: *big.NewFloat(helper.expectIMax),
	}
	return min, max
}
