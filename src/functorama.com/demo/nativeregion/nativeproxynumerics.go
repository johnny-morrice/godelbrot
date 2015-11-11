package nativeregion

import (
    "functorama.com/demo/region"
	"functorama.com/demo/nativesequence"
)

type NativeRegionProxy struct {
	*NativeRegionNumerics
	LocalRegion   nativeRegion
}

// Check we implement the interface
var _ region.RegionNumerics = NativeRegionProxy{}

func (proxy NativeRegionProxy) ClaimExtrinsics() {
	proxy.NativeRegionNumerics.Region = proxy.LocalRegion
}

type NativeSequenceProxy struct {
	*nativesequence.NativeSequenceNumerics
	LocalRegion   nativeRegion
}

func (proxy NativeSequenceProxy) ClaimExtrinsics() {
	base := proxy.NativeSequenceNumerics.NativeBaseNumerics
	rectangle := proxy.LocalRegion.rect(&base)
	proxy.NativeSequenceNumerics.SubImage(rectangle)
}