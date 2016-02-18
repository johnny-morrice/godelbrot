package nativesequence

import (
	"testing"
	"github.com/johnny-morrice/godelbrot/base"
	"github.com/johnny-morrice/godelbrot/nativebase"
)

func TestSequence(t *testing.T) {
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
	}
	app.PlaneMin = complex(0.0, 0.0)
	app.PlaneMax = complex(10.0, 10.0)
	numerics := Make(app)
	out := numerics.Sequence()

	const expectedCount = 100
	actualArea := numerics.Area()

	if expectedCount != actualArea {
		t.Error("Expected area of", expectedCount,
			"but received", actualArea)
	}

	members := make([]base.PixelMember, actualArea)

	i := 0
	for point := range out {
		members[i] = point
		i++
	}
	actualCount := len(members)

	if expectedCount != actualCount {
		t.Error("Expected", expectedCount, "members but there were", actualCount)
	}
}
