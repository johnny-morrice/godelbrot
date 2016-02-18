package nativeregion

import (
    "github.com/johnny-morrice/godelbrot/base"
    "github.com/johnny-morrice/godelbrot/region"
    "github.com/johnny-morrice/godelbrot/nativebase"
)

type RenderApplication interface {
    nativebase.NativeCoordProvider
    region.RegionProvider
    base.RenderApplication
}