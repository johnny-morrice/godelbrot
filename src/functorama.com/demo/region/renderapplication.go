package region

import (
    "functorama.com/demo/base"
)

type RenderApplication interface {
    base.RenderApplication
    // Configuration for particular render strategies
    RegionConfig() RegionParameters
    Factory() RegionNumericsFactory
}

type RegionParameters struct {
    GlitchSamples      uint
    RegionCollapseSize uint
}