package region

import (
    "functorama.com/demo/base"
    "functorama.com/demo/draw"
)

type MockRenderApplication struct {
    base.MockRenderApplication
    draw.MockContextProvider

    TRegionConfig bool
    TRegionNumericsFactory bool

    RegConfig RegionConfig
    RegionFactory *MockFactory
}

func (mock *MockRenderApplication) RegionConfig() RegionConfig {
    mock.TRegionConfig = true
    return mock.RegConfig
}

func (mock *MockRenderApplication) RegionNumericsFactory() RegionNumericsFactory {
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