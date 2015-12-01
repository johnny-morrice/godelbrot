package nativesharedregion

import (
	"functorama.com/demo/sharedregion"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
	"functorama.com/demo/nativebase"
)

type MockRenderApplication struct {
    nativebase.MockNativeCoordProvider
    region.MockRegionProvider
    base.MockRenderApplication
    sharedregion.MockSharedRegionProvider
}