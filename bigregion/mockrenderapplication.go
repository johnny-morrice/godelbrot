package bigregion

import (
    "github.com/johnny-morrice/godelbrot/bigbase"
    "github.com/johnny-morrice/godelbrot/region"
    "github.com/johnny-morrice/godelbrot/base"
)

type MockRenderApplication struct {
    bigbase.MockBigCoordProvider
    region.MockRegionProvider
    base.MockRenderApplication
}