package sharedregion

import (
    "functorama.com/demo/region"
)

type ConcurrentRegionParameters struct {
    BufferSize uint
    RenderJobs uint
}

type SharedRegionRenderApplication interface {
    RegionRenderApplication
    ConcurrentConfig() ConcurrentRegionParameters
    SharedRegionNumericsFactory() SharedRegionNumericsFactory
}