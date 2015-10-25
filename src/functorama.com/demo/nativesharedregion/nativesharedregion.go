package nativesharedregion

import (
    "functorama.com/demo/sharedregion"
    "functorama.com/demo/nativeregion"
)

type NativeSharedRegionNumerics struct {
    nativeregion.NativeRegionNumericsProxy
    prototypes []nativeregion.NativeRegionNumerics
}

func CreateNativeSharedRegionNumerics(numerics *nativeregion.NativeRegionNumerics, jobs uint) NativeSharedRegionNumerics {
    shared := NativeSharedRegionNumerics{
        prototypes: make([]nativeregion.NativeRegionNumerics, jobs),
    }
    for i := uint(0); i < jobs; i++ {
        shared.prototypes[i] = *numerics
    }
    initLocal := &shared.prototypes[0]
    shared.NativeRegionNumericsProxy = nativeregion.NativeRegionNumericsProxy{
        NativeRegionNumerics: initLocal,
        Region: initLocal.Region,
    }

    return shared
}

func (shared NativeSharedRegionNumerics) GrabThreadPrototype(threadId uint) {
    shared.NativeRegionNumericsProxy.NativeRegionNumerics = &shared.prototypes[threadId]
}

func (shared NativeSharedRegionNumerics) SharedChildren() []sharedregion.SharedRegionNumerics {
    localRegions := shared.NativeChildRegions()
    sharedChildren := make([]sharedregion.SharedRegionNumerics, len(localRegions))
    myCore := shared.NativeRegionNumericsProxy.NativeRegionNumerics
    for i, child := range localRegions {
        sharedChildren[i] = NativeSharedRegionNumerics{
            NativeRegionNumericsProxy: nativeregion.NativeRegionNumericsProxy{
                Region: child,
                NativeRegionNumerics: myCore,
            },
            prototypes: shared.prototypes,
        }
    }
    return sharedChildren
}