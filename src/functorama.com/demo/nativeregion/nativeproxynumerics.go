package nativeregion

import (
	"functorama.com/demo/nativesequence"
)

type NativeRegionProxy struct {
	*NativeRegionNumerics
	LocalRegion   NativeRegion
}

func (proxy NativeRegionProxy) ClaimExtrinsics() {
	proxy.NativeRegionNumerics.Region = proxy.LocalRegion
}

type NativeSequenceProxy struct {
	*nativesequence.NativeSequenceNumerics
	LocalRegion   NativeRegion
}

func (proxy NativeSequenceProxy) ClaimExtrinsics() {
	base := proxy.NativeSequenceNumerics.NativeBaseNumerics
	rectangle := proxy.LocalRegion.rect(&base)
	proxy.NativeSequenceNumerics.SubImage(rectangle)
}