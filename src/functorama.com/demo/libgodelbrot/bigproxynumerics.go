package libgodelbrot

type BigRegionNumericsProxy struct {
    Region BigRegion
    Numerics *BigRegionNumerics
}

func (proxy BigRegionNumericsProxy) Initialize() {
    proxy.Numerics.region = Region
}

type BigSequenceNumericsProxy struct {
    Region BigRegion
    Numerics *BigSequentialNumerics
}

func (proxy BigSequentialNumerics) Initialize() {
    proxy.Numerics.SubImage(proxy.Region.Rect())
}

