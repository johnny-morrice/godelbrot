package nativeregion

type NativeRegionNumericsProxy struct {
	region   nativeRegion
	numerics *NativeRegionNumerics
}

func (proxy NativeRegionNumericsProxy) ClaimExtrinsics() {
	proxy.numerics.region = proxy.region
}

type NativeSequenceNumericsProxy struct {
	region   nativeRegion
	numerics *NativeSequentialNumerics
}

func (proxy NativeSequentialNumerics) ClaimExtrinsics() {
	proxy.numerics.SubImage(proxy.region.Rect())
}
