package nativeregion

import (
    "github.com/johnny-morrice/godelbrot/internal/nativebase"
    "github.com/johnny-morrice/godelbrot/internal/region"
    "github.com/johnny-morrice/godelbrot/internal/base"
)

type MockRenderApplication struct {
    nativebase.MockNativeCoordProvider
    region.MockRegionProvider
    base.MockRenderApplication
}

var _ RenderApplication = (*MockRenderApplication)(nil)