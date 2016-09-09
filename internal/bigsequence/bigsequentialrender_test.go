package bigsequence

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/bigbase"
	"testing"
)

func TestBigMandelbrotSequence(t *testing.T) {
	const prec = 53
	const iterLimit = 10

	app := &bigbase.MockRenderApplication{
		MockRenderApplication: base.MockRenderApplication{
			Base: base.BaseConfig{
				DivergeLimit: 4.0,
			},
			PictureWidth:  10,
			PictureHeight: 10,
		},
	}
	app.UserMin = bigbase.MakeBigComplex(0.0, 0.0, prec)
	app.UserMax = bigbase.MakeBigComplex(10.0, 10.0, prec)
	numerics := Make(app)
	out := numerics.Sequence()

	const expectedCount = 100
	actualCount := len(out)

	if expectedCount != actualCount {
		t.Error("Expected", expectedCount, "members but there were", actualCount)
	}
}
