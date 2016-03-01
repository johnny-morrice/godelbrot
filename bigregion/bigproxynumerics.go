package bigregion

import (
    "github.com/johnny-morrice/godelbrot/sequence"
    "github.com/johnny-morrice/godelbrot/region"
    "github.com/johnny-morrice/godelbrot/bigsequence"
    "github.com/johnny-morrice/godelbrot/bigbase"
)

type BigRegionNumericsProxy struct {
    *BigRegionNumerics
	LocalRegion   bigRegion
}

var _ region.RegionNumerics = BigRegionNumericsProxy{}

func (proxy BigRegionNumericsProxy) ClaimExtrinsics() {
	proxy.BigRegionNumerics.Region = proxy.LocalRegion
}

func (proxy BigRegionNumericsProxy) Extrinsically(f func()) {
    old := proxy.BigRegionNumerics.Region
    proxy.ClaimExtrinsics()
    f()
    proxy.BigRegionNumerics.Region = old
}

type BigSequenceNumericsProxy struct {
    *bigsequence.BigSequenceNumerics
	LocalRegion   bigRegion
}

var _ sequence.SequenceNumerics = BigSequenceNumericsProxy{}

// TODO remove method.  Use Extrinsically instead.
func (proxy BigSequenceNumericsProxy) ClaimExtrinsics() {
    base := proxy.BigSequenceNumerics.BigBaseNumerics
    rectangle := proxy.LocalRegion.rect(&base)
    proxy.BigSequenceNumerics.SubImage(rectangle)
}

func (proxy BigSequenceNumericsProxy) Extrinsically(f func()) {
    cmin := bigbase.BigComplex{proxy.RealMin, proxy.ImagMin}
    cmax := bigbase.BigComplex{proxy.RealMax, proxy.ImagMax}

    proxy.ClaimExtrinsics()
    f()
    proxy.RealMin = cmin.R
    proxy.ImagMin = cmin.I
    proxy.RealMax = cmax.R
    proxy.ImagMax = cmax.I

    // Should we cache this somewhere?  New object?
    proxy.RestorePicBounds()
}