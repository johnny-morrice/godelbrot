package libgodelbrot

import (
    "functorama.com/demo/nativesharedregion"
)

type nativeSharedRegionFacade struct {
    *baseFacade
    *sharedRegionProvider
    *nativeCoords
}

var _ nativesharedregion.RenderApplication = (*nativeSharedRegionFacade)(nil)

func makeNativeSharedRegionFacade(desc *Info, baseApp *baseFacade, provider *sharedRegionProvider) *nativeSharedRegionFacade {
    facade := &nativeSharedRegionFacade{}
    facade.baseFacade = baseApp
    facade.sharedRegionProvider = provider
    facade.nativeCoords = makeNativeCoords(desc)
    return facade
}