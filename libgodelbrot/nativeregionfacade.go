package libgodelbrot

import (
    "functorama.com/demo/nativeregion"
)

type nativeRegionFacade struct {
    *baseFacade
    *regionProvider
    *nativeCoords
}

var _ nativeregion.RenderApplication = (*nativeRegionFacade)(nil)

func makeNativeRegionFacade(desc *Info, baseApp *baseFacade, regionDesc *regionProvider) *nativeRegionFacade {
    facade := &nativeRegionFacade{}
    facade.baseFacade = baseApp
    facade.regionProvider = regionDesc
    facade.nativeCoords = makeNativeCoords(desc)
    return facade
}