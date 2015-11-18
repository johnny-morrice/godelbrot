package nativeregion

import (
    "functorama.com/demo/base"
    "functorama.com/demo/region"
    "functorama.com/demo/nativebase"
)

type RenderApplication interface {
    nativebase.NativeCoordProvider
    region.RegionProvider
    base.RenderApplication
}