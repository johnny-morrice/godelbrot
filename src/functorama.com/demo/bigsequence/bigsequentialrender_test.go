package libgodelbrot

import (
	"testing"
)

func TestBigMandelbrotSequence(t *testing.T) {
	base := BaseNumerics{
		PicXMin: 0,
		PicXMax: 10,
		PicYMin: 0,
		PicYMax: 10,
	}
	bigBase := BigBaseNumerics{
		BaseNumerics: base,
		RealMin:      CreateBigFloat(0.0, Prec64),
		RealMax:      CreateBigFloat(10.0, Prec64),
		ImagMin:      CreateBigFloat(0.0, Prec64),
		ImagMax:      CreateBigFloat(10.0, Prec64),
		DivergeLimit: CreateBigFloat(4.0, Prec64),
	}
	numerics := CreateBigSequentialNumerics(bigBase)
	numerics.MemberCaptureSequencer()
	members := numerics.CapturedMembers()

	expectedCount := 100
	actualCount := len(members)

	if expectedCount != actualCount {
		t.Error("Expected", expectedCount, "members but there were", actualCount)
	}
}
