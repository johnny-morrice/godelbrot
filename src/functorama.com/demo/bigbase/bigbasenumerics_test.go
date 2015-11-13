package bigbase

import (
	"testing"
	"functorama.com/demo/base"
)

// Three paths through CreateBigBaseNumerics
// 1. Aspect ratio is okay
// 2. Aspect ratio is too short
// 3. Aspect ratio is too thin

func TestCreateBigBaseNumerics(t *testing.T) {
	noChange := aspectRatioFixHelper{
		pictureW: 200,
		pictureH: 100,

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
		pictureW: 200,
		pictureH: 100,

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
		pictureW: 200,
		pictureH: 100,

		rMin: -1.0,
		rMax: 1.0,
		iMin: -0.1,
		iMax: 0.1,

		expectRMin: -1.0,
		expectRMax: 1.0,
		expectIMin: -0.1,
		expectIMax: 0.9,
	}

	tests := []aspectRatioFixHelper{noChange, fatter, taller}
	for _, test := range tests {
		testCreateBigBaseNumerics(t, test)
	}
}

func testCreateBigBaseNumerics(t *testing.T, helper aspectRatioFixHelper) {
	userMin, userMax := helper.planeCoords()
	expectMin, expectMax := helper.expectCoords()

	mock := &MockRenderApplication{
		MockRenderApplication: base.MockRenderApplication{
			PictureWidth:   helper.pictureW,
			PictureHeight:   helper.pictureH,
			Base: base.BaseConfig {
				FixAspect: true,
			},
		},
		UserMin: userMin,
		UserMax: userMax,
	}

	numerics := CreateBigBaseNumerics(mock)

	fixOkay := bigEq(&numerics.RealMin, expectMin.Real())
	fixOkay = fixOkay && bigEq(&numerics.ImagMin, expectMin.Imag())
	fixOkay = fixOkay && bigEq(&numerics.RealMax, expectMax.Real())
	fixOkay = fixOkay && bigEq(&numerics.ImagMax, expectMax.Imag())

	if !fixOkay {
		t.Error("Aspect ratio fix broken for helper:, ", helper,
			" received: ", numerics)
	}

	mockOkay := mock.TBigUserCoords && mock.TPictureDimensions
	mockOkay = mockOkay && mock.TBigUserCoords && mock.TBaseConfig

	if !mockOkay {
		t.Error("Expected method not called on mock", mock)
	}
}

func TestPlaneToPixel(t *testing.T) {
	numerics := BigBaseNumerics{
		RealMin:     CreateBigFloat(-1.0, testPrec),
		ImagMax:     CreateBigFloat(1.0, testPrec),
	}

	const imageWidth = 100
	const imageHeight = 100

	qA := BigComplex{
		CreateBigFloat(0.1, testPrec),
		CreateBigFloat(0.1, testPrec),
	}

	qB := BigComplex{
		CreateBigFloat(0.1, testPrec),
		CreateBigFloat(-0.1, testPrec),
	}

	qC := BigComplex{
		CreateBigFloat(-0.1, testPrec),
		CreateBigFloat(-0.1, testPrec),
	}

	qD := BigComplex{
		CreateBigFloat(-0.1, testPrec),
		CreateBigFloat(0.1, testPrec),
	}

	origin := BigComplex{
		CreateBigFloat(0.0, testPrec),
		CreateBigFloat(0.0, testPrec),
	}

	offset := BigComplex{
		CreateBigFloat(-1.0, testPrec),
		CreateBigFloat(1.0, testPrec),
	}

	const expectPixAx = 55
	const expectPixAY = 45

	const expectPixBx = 55
	const expectPixBy = 55

	const expectPixCx = 45
	const expectPixCy = 55

	const expectPixDx = 45
	const expectPixDy = 45

	const expectOx = 50
	const expectOy = 50

	const expectOffsetX = 0
	const expectOffsetY = 0

	points := []BigComplex{qA, qB, qC, qD, origin, offset}
	expectedXs := []int{
		expectPixAx,
		expectPixBx,
		expectPixCx,
		expectPixDx,
		expectOx,
		expectOffsetX,
	}
	expectedYs := []int{
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

const testPrec = 53