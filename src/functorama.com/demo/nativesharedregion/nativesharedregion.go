package nativesharedregion

import (
    "functorama.com/demo/sharedregion"
    "functorama.com/demo/nativesequence"
    "functorama.com/demo/nativeregion"
)

type NativeSharedRegion struct {
    nativeregion.NativeRegionProxy
    workerId uint16
    prototypes []*nativeregion.NativeRegionNumerics
    sequencePrototypes []*nativesequence.NativeSequenceNumerics
}

var _ sharedregion.SharedRegionNumerics = NativeSharedRegion{}

func Make(app RenderApplication) NativeSharedRegion {
    sharedConfig := app.SharedRegionConfig()
    // Add a job for the tracker
    jobs := sharedConfig.Jobs
    shared := NativeSharedRegion{
        prototypes: make([]*nativeregion.NativeRegionNumerics, jobs),
        sequencePrototypes: make([]*nativesequence.NativeSequenceNumerics, jobs),
    }
    regorig := nativeregion.Make(app)
    seqorig := *regorig.SequenceNumerics
    for i := uint16(0); i < jobs; i++ {
        // Copy numerics
        reg := &nativeregion.NativeRegionNumerics{}
        *reg = regorig
        shared.prototypes[i] = reg

        seq := &nativesequence.NativeSequenceNumerics{}
        *seq = seqorig
        shared.sequencePrototypes[i] = seq

        shared.prototypes[i].SequenceNumerics = seq
    }
    initLocal := shared.prototypes[0]
    shared.NativeRegionProxy = nativeregion.NativeRegionProxy{
        NativeRegionNumerics: initLocal,
        LocalRegion: initLocal.Region,
    }

    return shared
}

func (shared NativeSharedRegion) GrabWorkerPrototype(workerId uint16) {
    shared.NativeRegionNumerics = shared.prototypes[workerId]
    shared.workerId = workerId
}

func (shared NativeSharedRegion) SharedChildren() []sharedregion.SharedRegionNumerics {
    smallreg := shared.NativeChildRegions()
    children := make([]sharedregion.SharedRegionNumerics, len(smallreg))
    regnum := shared.NativeRegionNumerics
    for i, r := range smallreg {
        children[i] = NativeSharedRegion{
            NativeRegionProxy: nativeregion.NativeRegionProxy{
                LocalRegion: r,
                NativeRegionNumerics: regnum,
            },
            prototypes: shared.prototypes,
            sequencePrototypes: shared.sequencePrototypes,
        }
    }
    return children
}

func (shared NativeSharedRegion) SharedRegionSequence() sharedregion.SharedSequenceNumerics {
    return shared.NativeSharedSequence()
}

func (shared NativeSharedRegion) NativeSharedSequence() NativeSharedSequence {
    return NativeSharedSequence{
        NativeSequenceProxy: nativeregion.NativeSequenceProxy{
            NativeSequenceNumerics: shared.sequencePrototypes[shared.workerId],
            LocalRegion: shared.LocalRegion,
        },
        prototypes: shared.sequencePrototypes,
    }
}

type NativeSharedSequence struct {
    nativeregion.NativeSequenceProxy
    prototypes []*nativesequence.NativeSequenceNumerics
    workerId uint16
}

var _ sharedregion.SharedSequenceNumerics = NativeSharedSequence{}

func (shared NativeSharedSequence) GrabWorkerPrototype(workerId uint16) {
    shared.NativeSequenceNumerics = shared.prototypes[workerId]
    shared.workerId = workerId
}