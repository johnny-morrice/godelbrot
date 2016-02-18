package bigregion

import (
    "github.com/johnny-morrice/godelbrot/base"
    "github.com/johnny-morrice/godelbrot/bigbase"
    "github.com/johnny-morrice/godelbrot/region"
)

type RenderApplication interface {
    bigbase.BigCoordProvider
    region.RegionProvider
    base.RenderApplication
}