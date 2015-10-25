package nativesharedregion

import (
    "functorama.com/demo/sharedregion"
    "functorama.com/demo/nativeregion"
)

type NativeSharedRegionNumerics struct {
    *NativeRegionNumerics
    prototypes []NativeRegionNumerics
}

func CreateNativeSharedRegionNumerics(jobs uint) {
    shared := NativeSharedRegionNumerics{
        prototypes: make([]NativeRegionNumerics, jobs)
    }
    for var i uint = 0; i < jobs; i++ {
        prototypes[i] = NativeRegionNumerics{}
    }
}

func (shared NativeSharedRegionNumerics) GrabThreadPrototype(threadId int) {
    shared.NativeRegionNumerics = &shared.prototypes[threadId]
}