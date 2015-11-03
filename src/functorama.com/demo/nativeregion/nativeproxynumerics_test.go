package nativeregion

import (
	"testing"
	"functorama.com/demo/base"
	"functorama.com/demo/nativebase"
	"functorama.com/demo/nativesequence"
)

func TestNativeProxyRegionClaimExtrinsics(t *testing.T) {
	native := NativeRegionProxy{
		LocalRegion: NativeRegion{
			topLeft: nativeMandelbrotThunk{
				evaluated: true,
			},
		},
		NativeRegionNumerics: &NativeRegionNumerics{},
	}

	native.ClaimExtrinsics()

	if native.LocalRegion != native.NativeRegionNumerics.Region {
		t.Error("Expected extrinsics were not claimed")
	}
}

func TestNativeProxySequenceClaimExtrinsics(t *testing.T) {
	regMin := complex(-1, -1)
	regMax := complex(1, 1)

	planeMin := complex(-2, -2)
	planeMax := complex(2, 2)

	const picWidth uint = 100
	const picHeight uint = 100

	planeDim := planeMax - planeMin
	planeWidth := real(planeDim)
	planeHeight := imag(planeDim)

	rUnit, iUnit := nativebase.PixelUnits(picWidth, picHeight, planeWidth, planeHeight)

	numerics := nativesequence.NativeSequenceNumerics{
		NativeBaseNumerics: nativebase.NativeBaseNumerics{
			BaseNumerics: base.BaseNumerics{
				PicXMin: 0,
				PicYMin: 0,
				PicXMax: int(picWidth),
				PicYMax: int(picHeight),
			},
			RealMin: real(planeMin),
			RealMax: real(planeMax),
			ImagMin: imag(planeMin),
			ImagMax: imag(planeMax),
			Runit: rUnit,
			Iunit: iUnit,
		},
	}
	native := NativeSequenceProxy{
		LocalRegion:   createNativeRegion(regMin, regMax),
		NativeSequenceNumerics: &numerics,
	}

	native.ClaimExtrinsics()

	expect := base.BaseNumerics{
		PicXMin: 25,
		PicXMax: 75,
		PicYMin: 25,
		PicYMax: 75,
	}

	actual := native.NativeBaseNumerics.BaseNumerics

	if actual != expect {
		t.Error("Expected ", expect, "but received", actual)
	}
}
