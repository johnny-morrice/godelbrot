package region

import (
    "functorama.com/demo/base"
    "functorama.com/demo/draw"
)

type MockRenderApplication struct {
    base.MockRenderApplication
    draw.MockContextProvider
    MockRegionProvider
}

type MockRegionProvider struct {
    TRegionConfig bool
    TRegionNumericsFactory bool

    RegConfig RegionConfig
    RegionFactory *MockFactory
}

func (mock *MockRegionProvider) RegionConfig() RegionConfig {
    mock.TRegionConfig = true
    return mock.RegConfig
}

func (mock *MockRegionProvider) RegionNumericsFactory() RegionNumericsFactory {
    mock.TRegionNumericsFactory = true
    return mock.RegionFactory
}

type MockFactory struct {
    TBuild bool
    Numerics *MockNumerics
}

func (mock *MockFactory) Build() RegionNumerics {
    mock.TBuild = true
    return mock.Numerics
}