package sharedregion

import (
    "functorama.com/demo/region"
)

type MockFactory struct {
    TBuild bool

    Numerics *MockNumerics
}

func (factory *MockFactory) Build() SharedRegionNumerics {
    factory.TBuild = true
    return factory.Numerics
}

type MockSharedRegionProvider struct {
    TSharedRegionConfig bool
    TSharedRegionFactory bool

    SharedConfig SharedRegionConfig
    SharedFactory SharedRegionFactory
}

type MockRenderApplication struct {
    region.MockRenderApplication
    MockSharedRegionProvider
}


func (msrp *MockSharedRegionProvider) SharedRegionConfig() SharedRegionConfig {
    msrp.TSharedRegionConfig = true
    return msrp.SharedConfig
}

func (msrp *MockSharedRegionProvider) SharedRegionFactory() SharedRegionFactory {
    msrp.TSharedRegionFactory = true
    return msrp.SharedFactory
}
