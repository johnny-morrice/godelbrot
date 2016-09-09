package nativeregion

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/nativebase"
	"github.com/johnny-morrice/godelbrot/internal/region"
)

type RenderApplication interface {
	nativebase.NativeCoordProvider
	region.RegionProvider
	base.RenderApplication
}
