package region

import (
    "functorama.com/demo/base"
    "functorama.com/demo/draw"
)

type RenderApplication interface {
    base.RenderApplication
    draw.ContextProvider
    // Configuration for particular render strategies
    RegionConfig() RegionParameters
    Factory() RegionNumericsFactory
}

type RegionParameters struct {
    GlitchSamples      uint
    CollapseSize 	   uint
}