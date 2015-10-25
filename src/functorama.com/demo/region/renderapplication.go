package region

import (
    "functorama.com/demo/base"
)

type RegionRenderApplication interface {
    base.BaseRenderApplication
    // Configuration for particular render strategies
    RegionConfig() RegionParameters
}

type RegionParameters struct {
    GlitchSamples      uint
    RegionCollapseSize uint
}