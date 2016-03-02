package nativeregion

import (
    "github.com/johnny-morrice/godelbrot/nativebase"
    "github.com/johnny-morrice/godelbrot/region"
    "github.com/johnny-morrice/godelbrot/base"
)

type MockRenderApplication struct {
    nativebase.MockNativeCoordProvider
    region.MockRegionProvider
    base.MockRenderApplication
}

var _ RenderApplication = (*MockRenderApplication)(nil)