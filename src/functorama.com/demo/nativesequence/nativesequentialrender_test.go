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
	const iterateLimit = 10
	app := &nativebase.MockRenderApplication{
		MockRenderApplication: base.MockRenderApplication{
			PictureWidth: 10,
			PictureHeight: 10,
			Base: base.BaseConfig{DivergeLimit: 4.0, IterateLimit: iterateLimit},
		},
		PlaneMin: complex(0.0, 0.0),
		PlaneMax: complex(10.0, 10.0),
	}
	numerics := CreateNativeSequenceNumerics(app)
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
	const iterateLimit = 10
	context := draw.NewMockDrawingContext(iterateLimit)
	app := &nativebase.MockRenderApplication{
		MockRenderApplication: base.MockRenderApplication{
			PictureWidth: 10,
			PictureHeight: 10,
			Base: base.BaseConfig{DivergeLimit: 4.0, IterateLimit: iterateLimit},
		},
		PlaneMin: complex(0.0, 0.0),
		PlaneMax: complex(10.0, 10.0),
	}
	numerics := CreateNativeSequenceNumerics(app)
	numerics.ImageDrawSequencer(context)
	numerics.MandelbrotSequence(iterateLimit)

	if !(context.TPicture && context.TColors) {
		t.Error("Expected methods not called on drawing context")
	}
}