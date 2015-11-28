package sharedmodule

import (
    "functorama.com/demo/sharedregion"
    "functorama.com/demo/inheritsequence"
    "functorama.com/demo/inheritregion"
)

type SharedRegion struct {
    inheritregion.RegionProxy
    workerId uint16
    prototypes []*inheritregion.RegionNumerics
    sequencePrototypes []*inheritsequence.SequenceNumerics
}

var _ sharedregion.SharedRegionNumerics = SharedRegion{}

func MakeNumerics(app RenderApplication) SharedRegion {
    sharedConfig := app.SharedConfig()
    shared := SharedRegion{
        prototypes: make([]*inheritregion.RegionNumerics, sharedConfig.Jobs),
        sequencePrototypes: make([]*inheritsequence.SequenceNumerics, jobs),
    }
    for i := uint16(0); i < jobs; i++ {
        // Copy
        shared.prototypes[i] = bigregion.CreateRegionNumerics(app)
        shared.sequencePrototypes[i] = shared.prototypes[i].SequenceNumerics
    }
    initLocal := shared.prototypes[0]
    shared.RegionProxy = inheritregion.RegionProxy{
        RegionNumerics: initLocal,
        LocalRegion: initLocal.Region,
    }

    return shared
}

func (shared SharedRegion) GrabWorkerPrototype(workerId uint16) {
    shared.RegionNumerics = shared.prototypes[workerId]
    shared.workerId = workerId
}

func (shared SharedRegion) SharedChildren() []sharedregion.SharedRegionNumerics {
    localRegions := shared.ActualChildren()
    sharedChildren := make([]sharedregion.SharedRegionNumerics, len(localRegions))
    myCore := shared.RegionNumerics
    for i, child := range localRegions {
        sharedChildren[i] = SharedRegion{
            RegionProxy: inheritregion.RegionProxy{
                LocalRegion: child,
                RegionNumerics: myCore,
            },
            prototypes: shared.prototypes,
        }
    }
    return sharedChildren
}

func (shared SharedRegion) SharedRegionSequence() sharedregion.SharedSequenceNumerics {
    return shared.SharedSequence()
}

func (shared SharedRegion) SharedSequence() SharedSequence {
    return SharedSequence{
        SequenceProxy: inheritregion.SequenceProxy{
            SequenceNumerics: shared.sequencePrototypes[shared.workerId],
            LocalRegion: shared.Region,
        },
        prototypes: shared.sequencePrototypes,
    }
}

type SharedSequence struct {
    inheritregion.SequenceProxy
    prototypes []*inheritsequence.SequenceNumerics
    workerId uint16
}

var _ sharedregion.SharedSequenceNumerics = SharedSequence{}

func (shared SharedSequence) GrabWorkerPrototype(workerId uint16) {
    shared.SequenceNumerics = shared.prototypes[workerId]
    shared.workerId = workerId
}