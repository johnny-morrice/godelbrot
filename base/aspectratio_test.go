package base

import (
	"testing"
)

func TestPictureAspectRatio(t *testing.T) {
	expect := 16.0 / 9.0
	actual := PictureAspectRatio(1920, 1080)
	if actual != expect {
		t.Error("Expected aspect ratio", expect,
			"but was", actual)
	}
}

func TestAppPictureAspectRatio(t *testing.T) {
	expect := 16.0 / 9.0
	mock := MockRenderApplication{
		PictureWidth:  1920,
		PictureHeight: 1080,
	}
	actual := AppPictureAspectRatio(&mock)

	if actual != expect {
		t.Error("Expected app aspect ratio", expect,
			"but was", actual)
	}

	if !mock.TPictureDimensions {
		t.Error("Expected aspect ratio to be based on mock picture dimensions")
	}
}
