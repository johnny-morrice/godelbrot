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
    // TODO
}

type BigSequenceNumericsProxy struct {
    *bigsequence.BigSequenceNumerics
	LocalRegion   bigRegion
}

var _ sequence.SequenceNumerics = BigSequenceNumericsProxy{}

func (bsnp BigSequenceNumericsProxy) ClaimExtrinsics() {
	bsnp.BigSequenceNumerics.SubImage(bsnp.LocalRegion.rect(&bsnp.BigSequenceNumerics.BigBaseNumerics))
}

func (bsnp BigSequenceNumericsProxy) Extrinsically(f func()) {
    // TODO
}