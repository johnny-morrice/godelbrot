package sharedregion

import (
	"functorama.com/demo/draw"
	"functorama.com/demo/base"
    "functorama.com/demo/region"
)

type SharedRegionFactory interface {
	Build() SharedRegionNumerics
}

type SharedRegionConfig struct {
    BufferSize uint
    Jobs uint16
}

type SharedProvider interface {
    SharedRegionConfig() SharedRegionConfig
    SharedRegionFactory() SharedRegionFactory
}

type RenderApplication interface {
	base.RenderApplication
	draw.ContextProvider
    region.RegionProvider
    SharedProvider
}