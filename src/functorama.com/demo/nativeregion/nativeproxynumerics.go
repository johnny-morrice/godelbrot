package nativeregion

import (
	"functorama.com/demo/nativesequence"
)

type NativeRegionNumericsProxy struct {
	*NativeRegionNumerics
	region   nativeRegion
}

func (proxy NativeRegionNumericsProxy) ClaimExtrinsics() {
	proxy.NativeRegionNumerics.region = proxy.region
}

type NativeSequenceNumericsProxy struct {
	*nativesequence.NativeSequenceNumerics
	region   nativeRegion
}

func (proxy NativeSequenceNumericsProxy) ClaimExtrinsics() {
	base := proxy.NativeSequenceNumerics.NativeBaseNumerics
	rectangle := proxy.region.rect(&base)
	proxy.NativeSequenceNumerics.SubImage(rectangle)
}