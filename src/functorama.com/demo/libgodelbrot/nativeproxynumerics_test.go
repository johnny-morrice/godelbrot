package libgodelbrot

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
	native := NativeSequenceNumericsProxy{
		Numerics: &NativeSequenceNumerics{},
	}

	native.ClaimExtrinsincs()

	// This isn't great at the moment...
	// If we make it this far, then the test passes.
}
