package sharedregion

import (
    "functorama.com/demo/region"
)

type SharedRegionNumericsFactory interface {
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
    SharedRegionNumericsFactory() SharedRegionNumericsFactory
}