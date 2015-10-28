package sharedregion

import (
	"functorama.com/demo/sequence"
	"functorama.com/demo/region"
)

// SharedSequentialNumerics provides sequential (column-wise) rendering calculations for a threaded
// render strategy
type SharedSequenceNumerics interface {
    sequence.SequenceNumerics
    OpaqueThreadPrototype
}

// SharedRegionNumerics provides a RegionNumerics for threaded render stregies
type SharedRegionNumerics interface {
    region.RegionNumerics
    OpaqueThreadPrototype
    SharedChildren() []SharedRegionNumerics
    SharedRegionSequence() SharedSequenceNumerics
}