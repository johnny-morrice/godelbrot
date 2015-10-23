package libgodelbrot

import (
	"math/big"
	"testing"
)

// Three paths through CreateNativeBaseNumerics
// 1. Aspect ratio is okay
// 2. Aspect ratio is too short
// 3. Aspect ratio is too thin

func TestCreateNativeBaseNumerics(t *testing.T) {
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

	tests := []aspectRatioFixHelper{noChange, fatter, taller}
	for _, test := range tests {
		testCreateNativeBaseNumerics(test)
	}
}

func testCreateNativeBaseNumerics(helper aspectRatioFixHelper) {
	userMin, userMax := helper.planeCoords()
	expectMin, expectMan := helper.planeCoords()

	mock := mockRenderApplication{
		pictureW:   helper.pictureW,
		pictureH:   helper.pictureH,
		bigUserMin: userMin,
		bigUserMax: userMax,
		fixAspect:  true,
	}

	numerics := CreateNativeBaseNumerics(mock)

	fixOkay := numerics.realMin == real(userMin)
	fixOkay = fixOkay && numerics.imagMin == imag(userMin)
	fixOkay = fixOkay && numerics.realMax == real(userMax)
	fixOkay = fixOkay && numerics.imagMax == imag(userMax)

	if !fixOkay {
		t.Error("Aspect ratio fix broken for helper:, ", helper,
			" received: ", numerics)
	}

	mockOkay := mock.tNativeUserCoords && tPictureDimensions
	mockOkay = mockOkay && tLimits && tNativeUserCoords
	mockOkay = mockOkay && tFixAspect

	if !mockOkay {
		t.Error("Expected method not called on mock", mock)
	}
}

func TestPlaneToPixel(t *testing.T) {
	numerics := NativeBaseNumerics{
		realMin:     -1.0,
		imagMax:     1.0,
		imageWidth:  100,
		imageHeight: 100,
	}

	const qA = 0.1 + 0.1i

	const qB = 0.1 - 0.1i

	const qC = -0.1 - 0.1i

	const qD = -0.1 + 0.1i

	const origin complex128 = 0

	const offset = complex(-1.0, 1.0)

	const expectPixAx uint = 55
	const expectPixAY uint = 45

	const expectPixBx uint = 55
	const expectPixBy uint = 55

	const expectPixCx uint = 45
	const expectPixCy uint = 55

	const expectPixDx uint = 45
	const expectPixDy uint = 45

	const expectOx uint = 50
	const expectOy uint = 50

	const expectOffsetX uint = 0
	const expectOffsetY uint = 0

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
