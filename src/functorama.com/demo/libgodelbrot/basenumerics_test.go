package libgodelbrot

import (
	"testing"
)

func TestCreateBaseNumerics(t *testing.T) {
	mock := mockRenderApplication{}
	mock.pictureWidth = 10
	mock.pictureHeight = 20
	mock.iterLimit = 200

	numerics := CreateBaseNumerics(mock)

	mockOkay := mock.tLimits && mock.tPictureWidth && mock.tPictureHeight
	if !mockOkay {
		t.Error("Expected method not called on mock", mock)
	}

	okay := numerics.picXMin == 0
	okay = okay && numerics.picXMax == 10
	okay = okay && numerics.picYMin == 0
	okay = okay && numerics.picYMax == 20
	okay = okay && numerics.iterLimit == 200

	if !okay {
		t.Error("numerics had unexpected value", numerics)
	}
}
