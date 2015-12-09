package bigsharedregion

import (
    "functorama.com/demo/sharedregion"
    "functorama.com/demo/bigsequence"
    "functorama.com/demo/bigregion"
)

type BigSharedRegion struct {
    bigregion.BigRegionNumericsProxy
    workerId uint16
    prototypes []*bigregion.BigRegionNumerics
    sequencePrototypes []*bigsequence.BigSequenceNumerics
}

var _ sharedregion.SharedRegionNumerics = BigSharedRegion{}

func MakeNumerics(app RenderApplication) BigSharedRegion {
    sharedConfig := app.SharedRegionConfig()
    jobs := sharedConfig.Jobs
    shared := BigSharedRegion{
        prototypes: make([]*bigregion.BigRegionNumerics, jobs),
        sequencePrototypes: make([]*bigsequence.BigSequenceNumerics, jobs),
    }
    numerics := bigregion.CreateBigRegionNumerics(app)
    for i := uint16(0); i < jobs; i++ {
        // Copy numerics
        another := numerics
        shared.prototypes[i] = &another
        shared.sequencePrototypes[i] = shared.prototypes[i].SequenceNumerics
    }
    initLocal := shared.prototypes[0]
    shared.BigRegionNumericsProxy = bigregion.BigRegionNumericsProxy{
        BigRegionNumerics: initLocal,
        LocalRegion: initLocal.Region,
    }

    return shared
}

func (shared BigSharedRegion) GrabWorkerPrototype(workerId uint16) {
    shared.BigRegionNumerics = shared.prototypes[workerId]
    shared.workerId = workerId
}

func (shared BigSharedRegion) SharedChildren() []sharedregion.SharedRegionNumerics {
    localRegions := shared.BigChildRegions()
    sharedChildren := make([]sharedregion.SharedRegionNumerics, len(localRegions))
    myCore := shared.BigRegionNumerics
    for i, child := range localRegions {
        sharedChildren[i] = BigSharedRegion{
            BigRegionNumericsProxy: bigregion.BigRegionNumericsProxy{
                LocalRegion: child,
                BigRegionNumerics: myCore,
            },
            prototypes: shared.prototypes,
        }
    }
    return sharedChildren
}

func (shared BigSharedRegion) SharedRegionSequence() sharedregion.SharedSequenceNumerics {
    return shared.SharedSequence()
}

func (shared BigSharedRegion) SharedSequence() BigSharedSequence {
    return BigSharedSequence{
        BigSequenceNumericsProxy: bigregion.BigSequenceNumericsProxy{
            BigSequenceNumerics: shared.sequencePrototypes[shared.workerId],
            LocalRegion: shared.Region,
        },
        prototypes: shared.sequencePrototypes,
    }
}

type BigSharedSequence struct {
    bigregion.BigSequenceNumericsProxy
    prototypes []*bigsequence.BigSequenceNumerics
    workerId uint16
}

var _ sharedregion.SharedSequenceNumerics = BigSharedSequence{}

func (shared BigSharedSequence) GrabWorkerPrototype(workerId uint16) {
    shared.BigSequenceNumericsProxy.BigSequenceNumerics = shared.prototypes[workerId]
    shared.workerId = workerId
}