package libgodelbrot

import (
	"testing"
)

func TestBigProxyRegionClaimExtrinsics(t *testing.T) {
	big := BigRegionNumericsProxy{
		Region: BigRegion{
			topLeft: bigMandelbrotThunk{
				evaluated: true,
			},
		},
		Numerics: &BigRegionNumerics{},
	}

	big.ClaimExtrinsincs()

	if big.Region != big.Numerics.region {
		t.Error("Expected extrinsics were not claimed")
	}
}

func TestBigProxySequenceClaimExtrinsics(t *testing.T) {
	const prec uint = 53
	regMin := CreateBigComplex(-1.0-1.0, prec)
	regMax := CreateBigComplex(1.0, 1.0, prec)

	planeMin := CreateBigComplex(-2.0, -2.0, prec)
	planeMax := CreateBigComplex(2.0, 2.0, prec)

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
		Region:   createBigRegion(regMin, regMax),
		Numerics: &numerics,
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
