package bigregion

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/bigbase"
	"github.com/johnny-morrice/godelbrot/internal/region"
)

type RenderApplication interface {
	bigbase.BigCoordProvider
	region.RegionProvider
	base.RenderApplication
}
