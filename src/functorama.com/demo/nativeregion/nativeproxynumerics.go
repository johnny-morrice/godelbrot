package nativeregion

import (
	"functorama.com/demo/nativesequence"
)

type NativeRegionNumericsProxy struct {
	*NativeRegionNumerics
	Region   NativeRegion
}

func (proxy NativeRegionNumericsProxy) ClaimExtrinsics() {
	proxy.NativeRegionNumerics.Region = proxy.Region
}

type NativeSequenceNumericsProxy struct {
	*nativesequence.NativeSequenceNumerics
	Region   NativeRegion
}

func (proxy NativeSequenceNumericsProxy) ClaimExtrinsics() {
	base := proxy.NativeSequenceNumerics.NativeBaseNumerics
	rectangle := proxy.Region.rect(&base)
	proxy.NativeSequenceNumerics.SubImage(rectangle)
}