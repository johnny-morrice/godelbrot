package nativesequence

import (
	"testing"
	"functorama.com/demo/base"
	"functorama.com/demo/nativebase"
	"functorama.com/demo/draw"
)

func TestMemberCaptureSequence(t *testing.T) {
	if testing.Short() {
		panic("nativesequence testing impossible in short mode")
	}
	base := base.BaseNumerics{
		PicXMin: 0,
		PicXMax: 10,
		PicYMin: 0,
		PicYMax: 10,
	}
	native := nativebase.NativeBaseNumerics{
		BaseNumerics: base,
		RealMin:      0.0,
		RealMax:      10.0,
		ImagMin:      0.0,
		ImagMax:      10.0,
		SqrtDivergeLimit: 2.0,
	}
	const iterateLimit = 10
	numerics := CreateNativeSequenceNumerics(native)
	numerics.MemberCaptureSequencer()
	numerics.MandelbrotSequence(iterateLimit)
	members := numerics.CapturedMembers()

	const expectedCount = 100
	actualCount := len(members)

	if expectedCount != actualCount {
		t.Error("Expected", expectedCount, "members but there were", actualCount)
	}
}

func TestImageDrawSequence(t *testing.T) {
	if testing.Short() {
		panic("nativesequence testing impossible in short mode")
	}
	base := base.BaseNumerics{
		PicXMin: 0,
		PicXMax: 10,
		PicYMin: 0,
		PicYMax: 10,
	}
	native := nativebase.NativeBaseNumerics{
		BaseNumerics: base,
		RealMin:      0.0,
		RealMax:      10.0,
		ImagMin:      0.0,
		ImagMax:      10.0,
		SqrtDivergeLimit: 2.0,
	}
	const iterateLimit = 10
	context := draw.NewMockDrawingContext(iterateLimit)
	numerics := CreateNativeSequenceNumerics(native)
	numerics.ImageDrawSequencer(context)
	numerics.MandelbrotSequence(iterateLimit)

	if !(context.TPicture && context.TColors) {
		t.Error("Expected methods not called on drawing context")
	}
}