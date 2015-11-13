package bigbase

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
		R: CreateBigFloat(helper.rMin, testPrec),
		I: CreateBigFloat(helper.iMin, testPrec),
	}
	max := BigComplex{
		R: CreateBigFloat(helper.rMax, testPrec),
		I: CreateBigFloat(helper.iMax, testPrec),
	}
	return min, max
}

func (helper aspectRatioFixHelper) expectCoords() (BigComplex, BigComplex) {
	min := BigComplex{
		R: CreateBigFloat(helper.expectRMin, testPrec),
		I: CreateBigFloat(helper.expectIMin, testPrec),
	}
	max := BigComplex{
		R: CreateBigFloat(helper.expectRMax, testPrec),
		I: CreateBigFloat(helper.expectIMax, testPrec),
	}
	return min, max
}
