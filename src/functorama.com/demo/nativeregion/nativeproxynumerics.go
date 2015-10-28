package nativeregion

import (
	"functorama.com/demo/nativesequence"
)

type NativeRegionProxy struct {
	*NativeRegionNumerics
	Region   NativeRegion
}

func (proxy NativeRegionProxy) ClaimExtrinsics() {
	proxy.NativeRegionNumerics.Region = proxy.Region
}

type NativeSequenceProxy struct {
	*nativesequence.NativeSequenceNumerics
	Region   NativeRegion
}

func (proxy NativeSequenceProxy) ClaimExtrinsics() {
	base := proxy.NativeSequenceNumerics.NativeBaseNumerics
	rectangle := proxy.Region.rect(&base)
	proxy.NativeSequenceNumerics.SubImage(rectangle)
}