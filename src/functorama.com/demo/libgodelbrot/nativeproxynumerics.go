package libgodelbrot

type NativeRegionNumericsProxy struct {
    Region NativeRegion
    Numerics *NativeRegionNumerics
}

func (proxy NativeRegionNumericsProxy) Initialize() {
    proxy.Numerics.region = Region
}

type NativeSequenceNumericsProxy struct {
    Region NativeRegion
    Numerics *NativeSequentialNumerics
}

func (proxy NativeSequentialNumerics) Initialize() {
    proxy.Numerics.SubImage(proxy.Region.Rect())
}

