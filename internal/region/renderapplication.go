package region

import (
    "github.com/johnny-morrice/godelbrot/internal/base"
    "github.com/johnny-morrice/godelbrot/internal/draw"
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