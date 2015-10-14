package libgodelbrot

type NativeRegionNumericsProxy struct {
    Region NativeRegion
    Numerics *NativeRegionNumerics
}

func (proxy NativeRegionNumericsProxy) ClaimExtrinsics() {
    proxy.Numerics.region = Region
}

type NativeSequenceNumericsProxy struct {
    Region NativeRegion
    Numerics *NativeSequentialNumerics
}

func (proxy NativeSequentialNumerics) ClaimExtrinsics() {
    proxy.Numerics.SubImage(proxy.Region.Rect())
}

