package libgodelbrot

import (
	"testing"
)

func TestBigMandelbrotSequence(t *testing.T) {
	base := BaseNumerics{
		picXMin: 0,
		picXMax: 10,
		picYMin: 0,
		picYMax: 10,
	}
	bigBase := BigBaseNumerics{
		BaseNumerics: base,
		realMin:      CreateBigFloat(0.0, Prec64),
		realMax:      CreateBigFloat(10.0, Prec64),
		imagMin:      CreateBigFloat(0.0, Prec64),
		imagMax:      CreateBigFloat(10.0, Prec64),
		divergeLimit: CreateBigFloat(4.0, Prec64),
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
