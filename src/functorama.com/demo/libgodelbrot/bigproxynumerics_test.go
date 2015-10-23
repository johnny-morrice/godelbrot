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
	big := BigSequenceNumericsProxy{
		Region: BigRegion{
			topLeft: bigMandelbrotThunk{},
		},
		Numerics: &BigSequenceNumerics{},
	}

	big.ClaimExtrinsincs()

	// This isn't great at the moment...
	// If we make it this far, then the test passes.
}
