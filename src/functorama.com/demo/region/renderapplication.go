package region

import (
    "functorama.com/demo/base"
    "functorama.com/demo/draw"
)

type RegionNumericsFactory interface {
    Build() RegionNumerics
}

type RenderApplication interface {
    base.RenderApplication
    draw.ContextProvider
    // Configuration for particular render strategies
    RegionConfig() RegionConfig
    RegionNumericsFactory() RegionNumericsFactory
}

type RegionConfig struct {
    GlitchSamples      uint
    CollapseSize 	   uint
}