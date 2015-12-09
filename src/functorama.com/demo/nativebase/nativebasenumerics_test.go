package nativebase

import (
	"testing"
	"functorama.com/demo/base"
)

// Three paths through Make
// 1. Aspect ratio is okay
// 2. Aspect ratio is too short
// 3. Aspect ratio is too thin

func TestMake(t *testing.T) {
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
		expectIMin: -0.9,
		expectIMax: 0.1,
	}

	tests := []aspectRatioFixHelper{noChange, fatter, taller}
	for _, test := range tests {
		testMake(t, test)
	}
}

func testMake(t *testing.T, helper aspectRatioFixHelper) {
	userMin, userMax := helper.planeCoords()
	expectMin, expectMax := helper.expectCoords()

	mock := &MockRenderApplication{
		MockRenderApplication: base.MockRenderApplication{
			Base: base.BaseConfig{
				FixAspect:  true,
			},
			PictureWidth:   helper.pictureW,
			PictureHeight:   helper.pictureH,
		},
	}
	mock.PlaneMin = userMin
	mock.PlaneMax = userMax

	numerics := Make(mock)


	actualMin := complex(numerics.RealMin, numerics.ImagMin)
	actualMax := complex(numerics.RealMax, numerics.ImagMax)

	if !(expectMin == actualMin && expectMax == actualMax) {
		t.Error("Aspect ratio fix broken.",
			"Expected", expectMin, expectMax,
			"but received", actualMin, actualMax,
			"(user input was", userMin, userMax, ")")
	}

	mockOkay := mock.TNativeUserCoords && mock.TPictureDimensions
	mockOkay = mockOkay && mock.TBaseConfig

	if !mockOkay {
		t.Error("Expected method not called on mock", mock)
	}
}

func TestPlaneToPixel(t *testing.T) {
	const side = 100
	const width = side
	const height = side
	const planeSide = 2.0
	numerics := NativeBaseNumerics{
		RealMin:     -1.0,
		ImagMin:	 -1.0,
		ImagMax:     1.0,
		RealMax:     1.0,
	}
	numerics.ImageWidth(width)
	numerics.ImageHeight(height)
	uq := UnitQuery{width, height, planeSide, planeSide}
	numerics.Runit, numerics.Iunit = uq.PixelUnits()

	const qA = 0.1 + 0.1i

	const qB = 0.1 - 0.1i

	const qC = -0.1 - 0.1i

	const qD = -0.1 + 0.1i

	const origin complex128 = 0

	const offset = complex(-1.0, 1.0)

	const expectPixAx int = 55
	const expectPixAY int = 45

	const expectPixBx int = 55
	const expectPixBy int = 55

	const expectPixCx int = 45
	const expectPixCy int = 55

	const expectPixDx int = 45
	const expectPixDy int = 45

	const expectOx int = 50
	const expectOy int = 50

	const expectOffsetX int = 0
	const expectOffsetY int = 0

	points := []complex128{qA, qB, qC, qD, origin, offset}
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
