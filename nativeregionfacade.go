package godelbrot

import (
    "github.com/johnny-morrice/godelbrot/internal/nativeregion"
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