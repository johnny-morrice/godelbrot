package bigsharedregion

import (
	"functorama.com/demo/sharedregion"
	"functorama.com/demo/region"
	"functorama.com/demo/bigbase"
	"functorama.com/demo/base"
)

type RenderApplication interface {
	sharedregion.SharedProvider
    bigbase.BigCoordProvider
    region.RegionProvider
    base.RenderApplication
}