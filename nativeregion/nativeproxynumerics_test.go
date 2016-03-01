package nativeregion

import (
	"testing"
	"github.com/johnny-morrice/godelbrot/base"
	"github.com/johnny-morrice/godelbrot/nativebase"
	"github.com/johnny-morrice/godelbrot/nativesequence"
)

const sqrtDLimit = float64(2.0)

func TestNativeProxyRegionClaimExtrinsics(t *testing.T) {
	native := NativeRegionProxy{
		LocalRegion: nativeRegion{
			topLeft: nativebase.NativeEscapeValue{C: 1},
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

	uq := nativebase.UnitQuery{picWidth, picHeight, planeWidth, planeHeight}
	rUnit, iUnit := uq.PixelUnits()

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
		LocalRegion:   createNativeRegion(regMin, regMax, sqrtDLimit),
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
