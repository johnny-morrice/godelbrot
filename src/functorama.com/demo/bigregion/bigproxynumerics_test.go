package bigregion

import (
	"testing"
	"functorama.com/demo/bigbase"
)

func TestBigProxyRegionClaimExtrinsics(t *testing.T) {
	big := BigRegionNumericsProxy{
		BigRegionNumerics: &BigRegionNumerics{},
		LocalRegion: bigRegion{
			topLeft: bigMandelbrotThunk{
				evaluated: true,
			},
		},
	}

	big.ClaimExtrinsics()

	if big.LocalRegion != big.BigRegionNumerics.Region {
		t.Error("Expected extrinsics were not claimed")
	}
}

func TestBigProxySequenceClaimExtrinsics(t *testing.T) {
	const prec uint = 53
	regMin := bigbase.CreateBigComplex(-1.0, -1.0, prec)
	regMax := bigbase.CreateBigComplex(1.0, 1.0, prec)

	planeMin := bigbase.CreateBigComplex(-2.0, -2.0, prec)
	planeMax := bigbase.CreateBigComplex(2.0, 2.0, prec)

	numerics := BigSequenceNumerics{
		RealMin: planeMin.Real(),
		RealMax: planeMax.Real(),
		ImagMin: planeMin.Imag(),
		ImagMax: planeMax.Imag(),
		BigBaseNumerics: BigBaseNumerics{
			BaseNumerics: BaseNumerics{
				PicXMin: 0,
				PicYMin: 0,
				PicXMax: 100,
				PicYMax: 100,
			},
		},
	}
	big := BigSequenceNumericsProxy{
		BigSequenceNumerics: &numerics,
		LocalRegion:   createBigRegion(regMin, regMax),
	}

	big.ClaimExtrinsincs()

	expect := BaseNumerics{
		PicXMin: 25,
		PicXMax: 75,
		PicYMin: 25,
		PicYMax: 25,
	}

	actual := big.BigBaseNumerics.BaseNumerics

	if actual != expect {
		t.Error("Expected ", expect, "but received", actual)
	}
}
