package sharedregion

import (
    "functorama.com/demo/region"
)

type SharedRegionFactory interface {
	Build() SharedRegionNumerics
}

type SharedRegionConfig struct {
    BufferSize uint
    RenderJobs uint
    Jobs uint32
}

type RenderApplication interface {
    region.RenderApplication
    SharedRegionConfig() SharedRegionConfig
    SharedRegionFactory() SharedRegionFactory
}