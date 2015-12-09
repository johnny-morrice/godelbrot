package bigsequence

import (
	"testing"
	"functorama.com/demo/base"
	"functorama.com/demo/bigbase"
)

func TestBigMandelbrotSequence(t *testing.T) {
	const prec = 53
	const iterLimit = 10

	app := &bigbase.MockRenderApplication{
		MockRenderApplication: base.MockRenderApplication{
			Base: base.BaseConfig{
				DivergeLimit: 4.0,
			},
			PictureWidth: 10,
			PictureHeight: 10,
		},
	}
	app.UserMin = bigbase.CreateBigComplex(0.0, 0.0, prec)
	app.UserMax = bigbase.CreateBigComplex(10.0, 10.0, prec)
	numerics := CreateBigSequenceNumerics(app)
	out := numerics.Sequence(iterLimit)

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
