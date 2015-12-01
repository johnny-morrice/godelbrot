package nativesharedregion

import (
	"functorama.com/demo/sharedregion"
	"functorama.com/demo/region"
	"functorama.com/demo/nativebase"
	"functorama.com/demo/base"
)

type RenderApplication interface {
	sharedregion.SharedProvider
    nativebase.NativeCoordProvider
    region.RegionProvider
    base.RenderApplication
}