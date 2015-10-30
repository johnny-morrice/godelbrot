package sharedregion

import (
    "functorama.com/demo/region"
)

type MockRenderApplication struct {
    region.MockRenderApplication
    TSharedRegionConfig bool
    TSharedRegionFactory bool

    SharedConfig SharedRegionConfig
    SharedFactory SharedRegionFactory
}


func (mock *MockRenderApplication) SharedRegionConfig() SharedRegionConfig {
    mock.TSharedRegionConfig = true
    return mock.SharedConfig
}

func (mock *MockRenderApplication) SharedRegionFactory() SharedRegionFactory {
    mock.TSharedRegionFactory = true
    return mock.SharedFactory
}
