package bigregion

import (
    "functorama.com/demo/bigbase"
    "functorama.com/demo/region"
    "functorama.com/demo/base"
)

type MockRenderApplication struct {
    bigbase.MockBigCoordProvider
    region.MockRegionProvider
    base.MockRenderApplication
}