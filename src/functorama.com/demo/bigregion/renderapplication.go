package bigregion

import (
    "functorama.com/demo/base"
    "functorama.com/demo/bigbase"
    "functorama.com/demo/region"
)

type RenderApplication interface {
    bigbase.BigCoordProvider
    region.RegionProvider
    base.RenderApplication
}