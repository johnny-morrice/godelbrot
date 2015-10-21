package libgodelbrot

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