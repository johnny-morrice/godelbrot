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
    numerics := nativeregion.Make(app)
    for i := uint16(0); i < jobs; i++ {
        // Copy numerics
        another := new(nativeregion.NativeRegionNumerics)
        *another = numerics
        shared.prototypes[i] = another
        shared.sequencePrototypes[i] = shared.prototypes[i].SequenceNumerics
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
    localRegions := shared.NativeChildRegions()
    sharedChildren := make([]sharedregion.SharedRegionNumerics, len(localRegions))
    myCore := shared.NativeRegionNumerics
    for i, child := range localRegions {
        sharedChildren[i] = NativeSharedRegion{
            NativeRegionProxy: nativeregion.NativeRegionProxy{
                LocalRegion: child,
                NativeRegionNumerics: myCore,
            },
            prototypes: shared.prototypes,
        }
    }
    return sharedChildren
}

func (shared NativeSharedRegion) SharedRegionSequence() sharedregion.SharedSequenceNumerics {
    return shared.NativeSharedSequence()
}

func (shared NativeSharedRegion) NativeSharedSequence() NativeSharedSequence {
    return NativeSharedSequence{
        NativeSequenceProxy: nativeregion.NativeSequenceProxy{
            NativeSequenceNumerics: shared.sequencePrototypes[shared.workerId],
            LocalRegion: shared.Region,
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