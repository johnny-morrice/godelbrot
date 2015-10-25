package nativesequence

import (
	"testing"
)

func TestNativeMandelbrotSequence(t *testing.T) {
	base := BaseNumerics{
		PicXMin: 0,
		PicXMax: 10,
		PicYMin: 0,
		PicYMax: 10,
	}
	bigBase := NativeBaseNumerics{
		BaseNumerics: base,
		RealMin:      0.0,
		RealMax:      10.0,
		ImagMin:      0.0,
		ImagMax:      10.0,
		DivergeLimit: 4.0,
	}
	numerics := CreateNativeSequentialNumerics(nativeBase)
	numerics.MemberCaptureSequencer()
	members := numerics.CapturedMembers()

	const expectedCount = 100
	actualCount := len(members)

	if expectedCount != actualCount {
		t.Error("Expected", expectedCount, "members but there were", actualCount)
	}
}
