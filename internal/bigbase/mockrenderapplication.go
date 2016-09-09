package bigbase

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
)

type MockRenderApplication struct {
	base.MockRenderApplication
	MockBigCoordProvider
}

var _ RenderApplication = (*MockRenderApplication)(nil)

type MockBigCoordProvider struct {
	TBigUserCoords bool
	TPrecision     bool

	UserMin BigComplex
	UserMax BigComplex
	Prec    uint
}

func (mbcp *MockBigCoordProvider) Precision() uint {
	mbcp.TPrecision = true
	return mbcp.Prec
}

func (mbcp *MockBigCoordProvider) BigUserCoords() (*BigComplex, *BigComplex) {
	mbcp.TBigUserCoords = true
	return &mbcp.UserMin, &mbcp.UserMax
}
