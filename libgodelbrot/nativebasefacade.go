package libgodelbrot

import (
    "math/big"
    "functorama.com/demo/nativebase"
)

type nativeCoords struct {
    userMin complex128
    userMax complex128
}

var _ nativebase.NativeCoordProvider = (*nativeCoords)(nil)

func (coords *nativeCoords) NativeUserCoords() (complex128, complex128) {
    return coords.userMin, coords.userMax
}

func makeNativeCoords(desc *Info) *nativeCoords {
    coords := &nativeCoords{}
    bigNums := []*big.Float{&desc.RealMin, &desc.ImagMin, &desc.RealMax, &desc.ImagMax}
    native := make([]float64, len(bigNums))

    for i, heapNum := range bigNums {
        native[i], _ = heapNum.Float64()
    }

    coords.userMin = complex(native[0], native[1])
    coords.userMax = complex(native[2], native[3])
    return coords
}

type nativeBaseFacade struct {
    *baseFacade
    *nativeCoords
}

var _ nativebase.RenderApplication = (*nativeBaseFacade)(nil)

func makeNativeBaseFacade(desc *Info, baseApp *baseFacade) *nativeBaseFacade {
    facade := &nativeBaseFacade{}
    facade.baseFacade = baseApp
    facade.nativeCoords = makeNativeCoords(desc)
    return facade
}