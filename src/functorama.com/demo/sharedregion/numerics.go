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

func SharedSequenceCollapse(numerics SharedRegionNumerics, workerId uint16) []base.PixelMember {
    collapse := numerics.SharedRegionSequence()
    collapse.GrabWorkerPrototype(workerId)
    var points []base.PixelMember
    collapse.Extrinsically(func () {
        points = sequence.Capture(collapse)
    })
    return points
}