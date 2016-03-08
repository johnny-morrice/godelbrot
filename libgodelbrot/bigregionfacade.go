package libgodelbrot

import (
    "github.com/johnny-morrice/godelbrot/internal/bigregion"
)

type bigRegionFacade struct {
    *baseFacade
    *regionProvider
    *bigCoords
}

var _ bigregion.RenderApplication = (*bigRegionFacade)(nil)

func makeBigRegionFacade(desc *Info, baseApp *baseFacade, region *regionProvider) *bigRegionFacade {
    facade := &bigRegionFacade{}
    facade.baseFacade = baseApp
    facade.regionProvider = region
    facade.bigCoords = makeBigCoords(desc)
    return facade
}