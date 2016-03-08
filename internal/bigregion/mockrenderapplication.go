package bigregion

import (
    "github.com/johnny-morrice/godelbrot/internal/bigbase"
    "github.com/johnny-morrice/godelbrot/internal/region"
    "github.com/johnny-morrice/godelbrot/internal/base"
)

type MockRenderApplication struct {
    bigbase.MockBigCoordProvider
    region.MockRegionProvider
    base.MockRenderApplication
}