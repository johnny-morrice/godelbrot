package bigbase

import (
	"math/big"
)

type aspectRatioFixHelper struct {
	pictureW uint
	pictureH uint

	rMin big.Float
	rMax big.Float
	iMin big.Float
	iMax big.Float

	expectRMin big.Float
	expectRMax big.Float
	expectIMin big.Float
	expectIMax big.Float
}

func (helper aspectRatioFixHelper) bigPlaneCoords() (BigComplex, BigComplex) {
	min := BigComplex{
		R: helper.rMin,
		I: helper.iMin,
	}
	max := BigComplex{
		R: helper.rMax,
		I: helper.iMax,
	}
	return min, max
}

func (helper aspectRatioFixHelper) bigExpectCoords() (BigComplex, BigComplex) {
	min := BigComplex{
		R: helper.expectRMin,
		I: helper.expectIMin,
	}
	max := BigComplex{
		R: helper.expectRMax,
		I: helper.expectIMax,
	}
	return min, max
}
