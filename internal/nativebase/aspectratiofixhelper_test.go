package nativebase

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

func (helper aspectRatioFixHelper) planeCoords() (complex128, complex128) {
    min := complex(helper.rMin, helper.iMin)
    max := complex(helper.rMax, helper.iMax)
    return min, max
}

func (helper aspectRatioFixHelper) expectCoords() (complex128, complex128) {
    min := complex(helper.expectRMin, helper.expectIMin)
    max := complex(helper.expectRMax, helper.expectIMax)
    return min, max
}
