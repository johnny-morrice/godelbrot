package libgodelbrot

import (
    "functorama.com/demo/bigsharedregion"
)

type bigSharedRegionFacade struct {
    *baseFacade
    *sharedRegionProvider
    *bigCoords
}

var _ bigsharedregion.RenderApplication = (*bigSharedRegionFacade)(nil)

func makeBigSharedRegionFacade(desc *Info, baseApp *baseFacade, provider *sharedRegionProvider) *bigSharedRegionFacade {
    facade := &bigSharedRegionFacade{}
    facade.baseFacade = baseApp
    facade.sharedRegionProvider = provider
    facade.bigCoords = makeBigCoords(desc)
    return facade
}