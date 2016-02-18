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
		R: MakeBigFloat(helper.rMin, testPrec),
		I: MakeBigFloat(helper.iMin, testPrec),
	}
	max := BigComplex{
		R: MakeBigFloat(helper.rMax, testPrec),
		I: MakeBigFloat(helper.iMax, testPrec),
	}
	return min, max
}

func (helper aspectRatioFixHelper) expectCoords() (BigComplex, BigComplex) {
	min := BigComplex{
		R: MakeBigFloat(helper.expectRMin, testPrec),
		I: MakeBigFloat(helper.expectIMin, testPrec),
	}
	max := BigComplex{
		R: MakeBigFloat(helper.expectRMax, testPrec),
		I: MakeBigFloat(helper.expectIMax, testPrec),
	}
	return min, max
}
