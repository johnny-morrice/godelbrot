package sharedregion

import (
    "functorama.com/demo/base"
	"functorama.com/demo/sequence"
	"functorama.com/demo/region"
)

// Copy a prototypical object instance into the local thread
type OpaqueThreadPrototype interface {
    GrabThreadPrototype(threadId uint)
}

// SharedSequentialNumerics provides sequential (column-wise) rendering calculations for a threaded
// render strategy
type SharedSequenceNumerics interface {
    region.ProxySequence
    OpaqueThreadPrototype
}

// SharedRegionNumerics provides a RegionNumerics for threaded render stregies
type SharedRegionNumerics interface {
    region.RegionNumerics
    OpaqueThreadPrototype
    SharedChildren() []SharedRegionNumerics
    SharedRegionSequence() SharedSequenceNumerics
}

func SharedSequenceCollapse(numerics SharedRegionNumerics, threadId uint, iterateLimit uint8) []base.PixelMember {
    collapse := numerics.SharedRegionSequence()
    collapse.GrabThreadPrototype(threadId)
    collapse.ClaimExtrinsics()
    return sequence.Capture(collapse, iterateLimit)
}