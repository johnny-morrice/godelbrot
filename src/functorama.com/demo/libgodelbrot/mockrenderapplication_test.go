package libgodelbrot

import (
	"functorama.com/demo/base"
	"functorama.com/demo/sequence"
	"functorama.com/demo/region"
	"functorama.com/demo/sharedregion"
)

type MockRenderApplication struct {
	base.MockRenderApplication
	sequence.MockRenderApplication
	region.MockRenderApplication
	sharedregion.MockRenderApplication

	bigUserMin         BigComplex
	bigUserMax         BigComplex
	nativeUserMin      complex128
	nativeUserMax      complex128
}

func (mock *mockRenderApplication) Limits() (uint, float64) {
	mock.tLimits = true
	return mock.iterLimit, mock.divergeLimit
}

func (mock *mockRenderApplication) PictureAspect() float64 {
	mock.tPictureAspect = true
	return mock.pictureAspect
}

func (mock *mockRenderApplication) BigUserCoords() (BigComplex, BigComplex) {
	mock.tBigUserCoords = true
	return mock.bigUserMin, mock.bigUserMax
}

func (mock *mockRenderApplication) NativeUserCoords() (complex128, complex128) {
	mock.tNativeUserCoords = true
	return mock.nativeUserMin, mock.nativeUserMax
}

func (mock *mockRenderApplication) FixAspect() bool {
	mock.tFixAspect = true
	return mock.fixAspect
}

func (mock *mockRenderApplication) SequentialNumerics() SequentialNumerics {
	mock.tSequentialNumerics = true
	return mock.sequentialNumerics
}

func (mock *mockRenderApplication) RegionNumerics() RegionNumerics {
	mock.tRegionNumerics = true
	return mock.regionNumerics
}

func (mock *mockRenderApplication) RegionConfig() RegionParameters {
	mock.tRegionConfig = true
	return mock.regionConfig
}

func (mock *mockRenderApplication) ConcurrentConfig() ConcurrentRegionParameters {
	mock.tConcurrentConfig = true
	return mock.concurrentConfig
}

func (mock *mockRenderApplication) Draw() DrawingContext {
	mock.tDraw = true
	return mock.draw
}
