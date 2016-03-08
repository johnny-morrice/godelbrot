package nativesequence

import (
	"testing"
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/nativebase"
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
	actualCount := len(out)

	if expectedCount != actualCount {
		t.Error("Expected", expectedCount, "members but there were", actualCount)
	}
}
