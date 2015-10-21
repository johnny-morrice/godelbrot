package libgodelbrot

import (
    "testing"
    "math/big"
)

// Three paths through CreateBigBaseNumerics
// 1. Aspect ratio is okay
// 2. Aspect ratio is too short
// 3. Aspect ratio is too thin

func TestCreateBigBaseNumerics(t *testing.T) {
    noChange := aspectRatioFixHelper{
        imageW: 200,
        imageH: 100,

        rMin: -1.0,
        rMax: 1.0,
        iMin: -0.5,
        iMax: 0.5,

        expectRMin: -1.0,
        expectRMax: 1.0,
        expectIMin: -0.5,
        expectIMax: 0.5,
    }
    fatter := aspectRatioFixHelper{
        imageW: 200,
        imageH: 100,

        rMin: -0.5,
        rMax: 0.5,
        iMin: -0.5,
        iMax: 0.5,

        expectRMin: -0.5,
        expectRMax: 1.5,
        expectIMin: -0.5,
        expectIMax: 0.5,
    }
    taller := aspectRatioFixHelper{
        imageW: 200,
        imageH: 100,

        rMin: -1.0,
        rMax: 1.0,
        iMin: -0.1,
        iMax: 0.1,

        expectRMin: -1.0,
        expectRMax: 1.0,
        expectIMin: -0.1,
        expectIMax: 0.9,
    }

    tests := []aspectRatioFixHelper{noChange,fatter,taller}
    for _, test := range tests {
        testCreateBigBaseNumerics(test)
    }
}

func testCreateBigBaseNumerics(helper aspectRatioFixHelper) {
    userMin, userMax := helper.planeCoords()
    expectMin, expectMan := helper.planeCoords()

    mock := mockRenderApplication{
        pictureW: helper.pictureW,
        pictureH: helper.pictureH,
        bigUserMin: userMin,
        bigUserMax: userMax,
        fixAspect: true,
    }

    numerics := CreateBigBaseNumerics(mock)

    fixOkay := bigEq(numerics.realMin, userMin.Real())
    fixOkay = fixOkay && bigEq(numerics.imagMin, userMin.Imag())
    fixOkay = fixOkay && bigEq(numerics.realMax, userMax.Real())
    fixOkay = fixOkay && bigEq(numerics.imagMax, userMax.Imag())

    if !fixOkay {
        t.Error("Aspect ratio fix broken for helper:, ", helper, 
            " received: ", numerics)
    }

    mockOkay := mock.tBigUserCoords && tPictureDimensions 
    mockOkay = mockOkay && tLimits && tBigUserCoords
    mockOkay = mockOkay && tFixAspect

    if !mockOkay {
        t.Error("Expected method not called on mock", mock)
    }
}

func TestPlaneToPixel(t *testing.T) {
    numerics := BigBaseNumerics{
        realMin: CreateBigFloat(-1.0, Prec64), 
        imagMax: CreateBigFloat(1.0, Prec64),
        imageWidth: 100,
        imageHeight: 100,
    }

    qA := BigComplex{
        CreateBigFloat(0.1, Prec64), 
        CreateBigFloat(0.1, Prec64),
    }

    qB := BigComplex{
        CreateBigFloat(0.1, Prec64), 
        CreateBigFloat(-0.1, Prec64),
    }

    qC := BigComplex{
        CreateBigFloat(-0.1, Prec64), 
        CreateBigFloat(-0.1, Prec64),
    }

    qD := BigComplex{
        CreateBigFloat(-0.1, Prec64),
        CreateBigFloat(0.1, Prec64),
    }

    origin := BigComplex{
        CreateBigFloat(0.0, Prec64), 
        CreateBigFloat(0.0, Prec64),
    }

    offset := BigComplex{
        CreateBigFloat(-1.0, Prec64), 
        CreateBigFloat(1.0, Prec64),
    }

    var expectPixAx uint = 55
    var expectPixAY uint = 45

    var expectPixBx uint = 55
    var expectPixBy uint = 55

    var expectPixCx uint = 45
    var expectPixCy uint = 55

    var expectPixDx uint = 45
    var expectPixDy uint = 45

    var expectOx uint = 50
    var expectOy uint = 50

    var expectOffsetX uint = 0
    var expectOffsetY uint = 0

    points := []complex128{qA, qB, qC, qD, origin, offset}
    expectedXs := []uint{
        expectPixAx,
        expectPixBx,
        expectPixCx,
        expectPixDx,
        expectOx,
        expectOffsetX,
    }
    expectedYs := []uint{
        expectPixAY,
        expectPixBy,
        expectPixCy,
        expectPixDy,
        expectOy,
        expectOffsetY,
    }

    for i, point := range points {
        expectedX := expectedXs[i]
        expectedY := expectedYs[i]
        actualX, actualY := numerics.PlaneToPixel(point)
        if actualX != expectedX || actualY != expectedY {
            t.Error("Error on point", i, ":", point,
                " expected (", expectedX, ",", expectedY, ") but was",
                "(", actualX, ",", actualY, ")")
        }
    }
}

func (base *BigBaseNumerics) TestFastPixelPerfectPrecision(t *testing.T) {
    injective := bigPerfectPixelHelper(CreateBigFloat(1.0, Prec64))
    twentySeven := bigPerfectPixelHelper(CreateBigFloat(math.nextAfter32(0.0, 1.0), Prec64))
    fiftyThree := bigPerfectPixelHelper(CreateBigFloat(math.nextAfter(0.0, 1.0), Prec64))

    bases := []BigBaseNumerics{
        injective,
        twentySeven,
        fiftyThree,
    }

    expectations := []uint{
        7,
        27,
        53,
    }

    for i, base := range bases {
        expect := expectations[i]
        actual := base.FastPixelPerfectPrecision()

        if expect != actual {
            t.Error("At base", i, "expected precision ", expect, "but received", actual)
        }
        if !allBigPrecSet(base, actual){
            t.Error("Precision ", actual, "not set on base", i)
        }
    }
}

func TestSetPrec(t *testing.T) {
    prec := 98
    base := BigBaseNumerics{}
    base.SetPrec(prec)
    if !allBigPrecSet(base, prec) {
        t.Error("Precison not set on base")
    }
}

func allBigPrecSet(base BigBaseNumerics, prec uint) {
    okay := base.realMin.Prec() == prec
    okay = okay && base.realMax.Prec() == prec
    okay = okay && base.imagMin.Prec() == prec
    okay = okay && base.imagMax.Prec() == prec
    okay = okay && base.divergeLimit.Prec() == prec
    okay = okay && base.rUnit.Prec() == prec
    okay = okay && base.iUnit.Prec() == prec
    return okay
}

