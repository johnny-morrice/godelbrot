package libgodelbrot

import (
    "github.com/johnny-morrice/godelbrot/bigbase"
)

type bigCoords struct {
    userMin *bigbase.BigComplex
    userMax *bigbase.BigComplex
    precision uint
}

var _ bigbase.BigCoordProvider = (*bigCoords)(nil)

func makeBigCoords(desc *Info) *bigCoords {
    coords := &bigCoords{}
    coords.precision = desc.Precision
    coords.userMin = &bigbase.BigComplex{desc.RealMin, desc.ImagMin}
    coords.userMax = &bigbase.BigComplex{desc.RealMax, desc.ImagMax}
    return coords
}

func (coords *bigCoords) Precision() uint {
    return coords.precision
}

func (coords *bigCoords) BigUserCoords() (*bigbase.BigComplex, *bigbase.BigComplex) {
    return coords.userMin, coords.userMax
}

type bigBaseFacade struct {
    *baseFacade
    *bigCoords
}

var _ bigbase.RenderApplication = (*bigBaseFacade)(nil)

func makeBigBaseFacade(desc *Info, baseApp *baseFacade) *bigBaseFacade {
    app := &bigBaseFacade{}
    app.baseFacade = baseApp
    app.bigCoords = makeBigCoords(desc)
    return app
}