package libgodelbrot

type BigRegionNumericsProxy struct {
    Region BigRegion
    Numerics *BigRegionNumerics
}

func (proxy BigRegionNumericsProxy) ClaimExtrinsics() {
    proxy.Numerics.region = Region
}

type BigSequenceNumericsProxy struct {
    Region BigRegion
    Numerics *BigSequentialNumerics
}

func (proxy BigSequentialNumerics) ClaimExtrinsics() {
    proxy.Numerics.SubImage(proxy.Region.Rect())
}

