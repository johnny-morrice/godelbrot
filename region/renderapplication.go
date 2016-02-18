package region

import (
    "functorama.com/demo/base"
    "functorama.com/demo/draw"
)

type RegionNumericsFactory interface {
    Build() RegionNumerics
}

type RegionProvider interface {
    // Configuration for particular render strategies
    RegionConfig() RegionConfig
    RegionNumericsFactory() RegionNumericsFactory
}

type RenderApplication interface {
    base.RenderApplication
    draw.ContextProvider
    RegionProvider
}

type RegionConfig struct {
    Samples      uint
    CollapseSize 	   uint
}