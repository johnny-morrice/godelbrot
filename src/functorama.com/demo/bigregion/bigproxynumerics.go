package bigregion

import (
    "functorama.com/demo/sequence"
    "functorama.com/demo/region"
    "functorama.com/demo/bigsequence"
    "functorama.com/demo/bigbase"
)

type BigRegionNumericsProxy struct {
    *BigRegionNumerics
	LocalRegion   bigRegion
}

var _ region.RegionNumerics = BigRegionNumericsProxy{}

func (brnp BigRegionNumericsProxy) ClaimExtrinsics() {
	brnp.BigRegionNumerics.Region = brnp.LocalRegion
}

func (brnp BigRegionNumericsProxy) Extrinsically(f func()) {
    old := brnp.BigRegionNumerics.Region
    brnp.ClaimExtrinsics()
    f()
    brnp.BigRegionNumerics.Region = old
}

type BigSequenceNumericsProxy struct {
    *bigsequence.BigSequenceNumerics
	LocalRegion   bigRegion
}

var _ sequence.SequenceNumerics = BigSequenceNumericsProxy{}

// TODO remove method.  Use Extrinsically instead.
func (bsnp BigSequenceNumericsProxy) ClaimExtrinsics() {
    base := bsnp.BigSequenceNumerics.BigBaseNumerics
    rectangle := bsnp.LocalRegion.rect(&base)
    bsnp.BigSequenceNumerics.SubImage(rectangle)
}

func (bsnp BigSequenceNumericsProxy) Extrinsically(f func()) {
    cmin := bigbase.BigComplex{bsnp.RealMin, bsnp.ImagMin}
    cmax := bigbase.BigComplex{bsnp.RealMax, bsnp.ImagMax}

    bsnp.ClaimExtrinsics()
    f()
    bsnp.RealMin = cmin.R
    bsnp.ImagMin = cmin.I
    bsnp.RealMax = cmax.R
    bsnp.ImagMax = cmax.I

    // Should we cache this somewhere?  New object?
    bsnp.RestorePicBounds()
}