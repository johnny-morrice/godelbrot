package bigregion

import (
	"testing"
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/bigbase"
	"github.com/johnny-morrice/godelbrot/internal/bigsequence"
)

func TestBigProxyRegionClaimExtrinsics(t *testing.T) {
	const prec = 53
	const ilim = 255
	parent := bigbase.BigBaseNumerics{}
	parent.IterateLimit = ilim
	parent.Precision = prec
	parent.SqrtDivergeLimit = parent.MakeBigFloat(2.0)
	min := parent.MakeBigComplex(-1.0, -1.0)
	max := parent.MakeBigComplex(1.0, 1.0)

	big := BigRegionNumericsProxy{}
	big.BigRegionNumerics = &BigRegionNumerics{}
	big.LocalRegion = createBigRegion(parent, min, max)

	big.ClaimExtrinsics()

	if !regionEq(big.LocalRegion, big.BigRegionNumerics.Region) {
		t.Error("Expected", big.LocalRegion,
			"but received", big.BigRegionNumerics.Region)
	}
}

func TestBigProxySequenceClaimExtrinsics(t *testing.T) {
	const picW = 100
	const picH = 100

	app := &bigbase.MockRenderApplication{}
	app.PictureWidth = picW
	app.PictureHeight = picH
	app.UserMin = bigbase.MakeBigComplex(-2.0, -2.0, prec)
	app.UserMax = bigbase.MakeBigComplex(2.0, 2.0, prec)
	app.Prec = prec

	numerics := bigsequence.Make(app)

	min := bigbase.MakeBigComplex(-1.0, -1.0, prec)
	max := bigbase.MakeBigComplex(1.0, 1.0, prec)
	region := createBigRegion(numerics.BigBaseNumerics, min, max)

	bsnp := BigSequenceNumericsProxy{
		BigSequenceNumerics: &numerics,
		LocalRegion:   region,
	}

	bsnp.ClaimExtrinsics()

	expect := base.BaseNumerics{
		WholeWidth: picW,
		WholeHeight: picH,
		PicXMin: 25,
		PicXMax: 75,
		PicYMin: 25,
		PicYMax: 75,
	}

	actual := bsnp.BigBaseNumerics.BaseNumerics

	if actual != expect {
		t.Error("Expected ", expect, "but received", actual)
	}
}

// regionEq returns true when its inputs have the same locations in space
func regionEq(areg, breg bigRegion) bool {
	aths := areg.points()
	bths := breg.points()

	for i, a := range(aths) {
		b := bths[i]
		if !bigbase.BigComplexEq(a.C, b.C) {
			return false
		}
	}

	return true
}