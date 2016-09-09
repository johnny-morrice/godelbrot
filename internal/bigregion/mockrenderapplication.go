package bigregion

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/bigbase"
	"github.com/johnny-morrice/godelbrot/internal/region"
)

type MockRenderApplication struct {
	bigbase.MockBigCoordProvider
	region.MockRegionProvider
	base.MockRenderApplication
}
