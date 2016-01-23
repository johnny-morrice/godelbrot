package bigregion

import (
    "functorama.com/demo/sequence"
    "functorama.com/demo/region"
    "functorama.com/demo/bigsequence"
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

func (bsnp BigSequenceNumericsProxy) ClaimExtrinsics() {
    base := bsnp.BigSequenceNumerics.BigBaseNumerics
    rectangle := bsnp.LocalRegion.rect(&base)
    bsnp.BigSequenceNumerics.SubImage(rectangle)
}

func (bsnp BigSequenceNumericsProxy) Extrinsically(f func()) {
    // TODO
}