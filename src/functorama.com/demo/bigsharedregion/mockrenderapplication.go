package bigsharedregion

import (
	"functorama.com/demo/sharedregion"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
	"functorama.com/demo/bigbase"
)

type MockRenderApplication struct {
    bigbase.MockBigCoordProvider
    region.MockRegionProvider
    base.MockRenderApplication
    sharedregion.MockSharedRegionProvider
}