package bigregion

import (
	"testing"
	"github.com/johnny-morrice/godelbrot/base"
	"github.com/johnny-morrice/godelbrot/bigbase"
	"github.com/johnny-morrice/godelbrot/bigsequence"
)

func TestBigProxyRegionClaimExtrinsics(t *testing.T) {
	const prec = 53
	min := bigbase.MakeBigComplex(-1.0, -1.0, prec)
	max := bigbase.MakeBigComplex(1.0, 1.0, prec)

	big := BigRegionNumericsProxy{}
	big.BigRegionNumerics = &BigRegionNumerics{}
	big.LocalRegion = createBigRegion(min, max)

	big.ClaimExtrinsics()

	if !regionEq(big.LocalRegion, big.BigRegionNumerics.Region) {
		t.Error("Expected", big.LocalRegion,
			"but received", big.BigRegionNumerics.Region)
	}
}

func TestBigProxySequenceClaimExtrinsics(t *testing.T) {
	const picW = 100
	const picH = 100

	min := bigbase.MakeBigComplex(-1.0, -1.0, prec)
	max := bigbase.MakeBigComplex(1.0, 1.0, prec)
	region := createBigRegion(min, max)

	app := &bigbase.MockRenderApplication{}
	app.PictureWidth = picW
	app.PictureHeight = picH
	app.UserMin = bigbase.MakeBigComplex(-2.0, -2.0, prec)
	app.UserMax = bigbase.MakeBigComplex(2.0, 2.0, prec)

	numerics := bigsequence.Make(app)

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