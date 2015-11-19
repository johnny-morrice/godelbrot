package bigregion

import (
	"testing"
	"functorama.com/demo/base"
	"functorama.com/demo/bigbase"
	"functorama.com/demo/bigsequence"
)

func TestBigProxyRegionClaimExtrinsics(t *testing.T) {
	big := BigRegionNumericsProxy{}
	big.BigRegionNumerics = &BigRegionNumerics{}
	big.LocalRegion = bigRegion{
		topLeft: bigMandelbrotThunk{
			evaluated: true,
		},
	}

	big.ClaimExtrinsics()

	if !regionEq(big.LocalRegion, big.BigRegionNumerics.Region) {
		t.Error("Expected", big.LocalRegion,
			"but received", big.BigRegionNumerics.Region)
	}
}

func TestBigProxySequenceClaimExtrinsics(t *testing.T) {
	min := bigbase.CreateBigComplex(-1.0, -1.0, prec)
	max := bigbase.CreateBigComplex(1.0, 1.0, prec)
	region := createBigRegion(min, max)

	app := &bigbase.MockRenderApplication{}
	app.PictureWidth = 100
	app.PictureHeight = 100
	app.UserMin = bigbase.CreateBigComplex(-2.0, -2.0, prec)
	app.UserMax = bigbase.CreateBigComplex(2.0, 2.0, prec)

	numerics := bigsequence.CreateBigSequenceNumerics(app)

	bsnp := BigSequenceNumericsProxy{
		BigSequenceNumerics: &numerics,
		LocalRegion:   region,
	}

	bsnp.ClaimExtrinsics()

	expect := base.BaseNumerics{
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

// Regions are considered equal on basis of cStore value
func regionEq(areg, breg bigRegion) bool {
	aths := areg.thunks()
	bths := breg.thunks()

	for i, a := range(aths) {
		b := bths[i]
		if !bigbase.BigComplexEq(&a.cStore, &b.cStore) {
			return false
		}
	}

	return true
}