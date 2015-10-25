package nativeregion

import (
	"testing"
)

func TestNativeProxyRegionClaimExtrinsics(t *testing.T) {
	native := NativeRegionNumericsProxy{
		Region: NativeRegion{
			topLeft: nativeMandelbrotThunk{
				evaluated: true,
			},
		},
		Numerics: &NativeRegionNumerics{},
	}

	native.ClaimExtrinsincs()

	if native.Region != native.Numerics.region {
		t.Error("Expected extrinsics were not claimed")
	}
}

func TestNativeProxySequenceClaimExtrinsics(t *testing.T) {
	const prec uint = 53
	regMin := complex(-1, -1)
	regMax := complex(1, 1)

	planeMin := complex(-2, -2)
	planeMax := complex(2, 2)

	numerics := NativeSequenceNumerics{
		RealMin: real(planeMin),
		RealMax: real(planeMax),
		ImagMin: imag(planeMin),
		ImagMax: imag(planeMax),
		NativeBaseNumerics: NativeBaseNumerics{
			BaseNumerics: BaseNumerics{
				PicXMin: 0,
				PicYMin: 0,
				PicXMax: 100,
				PicYMax: 100,
			},
		},
	}
	native := NativeSequenceNumericsProxy{
		Region:   createNativeRegion(regMin, regMax),
		Numerics: &numerics,
	}

	native.ClaimExtrinsincs()

	expect := BaseNumerics{
		PicXMin: 25,
		PicXMax: 75,
		PicYMin: 25,
		PicYMax: 25,
	}

	actual := native.NativeBaseNumerics.BaseNumerics

	if actual != expect {
		t.Error("Expected ", expect, "but received", actual)
	}
}
