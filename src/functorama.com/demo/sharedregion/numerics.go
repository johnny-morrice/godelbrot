package sharedregion

import (
    "functorama.com/demo/base"
    "functorama.com/demo/sequence"
	"functorama.com/demo/region"
)

// Copy a prototypical object instance into the local thread
type OpaqueWorkerPrototype interface {
    GrabWorkerPrototype(workerId uint16)
}

// SharedSequentialNumerics provides sequential (column-wise) rendering calculations for a threaded
// render strategy
type SharedSequenceNumerics interface {
    region.ProxySequence
    OpaqueWorkerPrototype
}

// SharedRegionNumerics provides a RegionNumerics for threaded render stregies
type SharedRegionNumerics interface {
    region.RegionNumerics
    OpaqueWorkerPrototype
    SharedChildren() []SharedRegionNumerics
    SharedRegionSequence() SharedSequenceNumerics
}

func SharedSequenceCollapse(reg SharedRegionNumerics, wid uint16) []base.PixelMember {
    reg.GrabWorkerPrototype(wid)
    reg.ClaimExtrinsics()
    seq := reg.SharedRegionSequence()
    seq.GrabWorkerPrototype(wid)
    var px []base.PixelMember
    seq.Extrinsically(func () {
        px = sequence.Capture(seq)
    })
    return px
}