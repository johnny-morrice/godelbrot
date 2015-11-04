package nativesharedregion

import (
    "functorama.com/demo/sharedregion"
    "functorama.com/demo/nativesequence"
    "functorama.com/demo/nativeregion"
)

type NativeSharedRegion struct {
    nativeregion.NativeRegionProxy
    threadId uint
    prototypes []nativeregion.NativeRegionNumerics
    sequencePrototypes []nativesequence.NativeSequenceNumerics
}

func CreateNativeSharedRegion(numerics *nativeregion.NativeRegionNumerics, jobs uint) NativeSharedRegion {
    shared := NativeSharedRegion{
        prototypes: make([]nativeregion.NativeRegionNumerics, jobs),
        sequencePrototypes: make([]nativesequence.NativeSequenceNumerics, jobs),
    }
    for i := uint(0); i < jobs; i++ {
        shared.prototypes[i] = *numerics
        shared.sequencePrototypes[i] = *shared.prototypes[i].SequenceNumerics
    }
    initLocal := &shared.prototypes[0]
    shared.NativeRegionProxy = nativeregion.NativeRegionProxy{
        NativeRegionNumerics: initLocal,
        LocalRegion: initLocal.Region,
    }

    return shared
}

func (shared NativeSharedRegion) GrabThreadPrototype(threadId uint) {
    shared.NativeRegionNumerics = &shared.prototypes[threadId]
    shared.threadId = threadId
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
            NativeSequenceNumerics: &shared.sequencePrototypes[shared.threadId],
            LocalRegion: shared.Region,
        },
        prototypes: shared.sequencePrototypes,
    }
}

type NativeSharedSequence struct {
    nativeregion.NativeSequenceProxy
    prototypes []nativesequence.NativeSequenceNumerics
    threadId uint
}

func (shared NativeSharedSequence) GrabThreadPrototype(threadId uint) {
    shared.NativeSequenceNumerics = &shared.prototypes[threadId]
    shared.threadId = threadId
}