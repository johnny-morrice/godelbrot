package base

import (
	"testing"
)

func TestCreateBaseNumerics(t *testing.T) {
	mock := &MockRenderApplication{}
	mock.PictureWidth = 10
	mock.PictureHeight = 20

	numerics := Make(mock)

	mockOkay := mock.TPictureDimensions
	if !mockOkay {
		t.Error("Expected method not called on mock", mock)
	}

	okay := numerics.PicXMin == 0
	okay = okay && numerics.PicXMax == 10
	okay = okay && numerics.PicYMin == 0
	okay = okay && numerics.PicYMax == 20

	if !okay {
		t.Error("numerics had unexpected value", numerics)
	}
}
