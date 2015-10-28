package base

type mockRenderApplication struct {
	tLimits             bool
	tPictureDimensions  bool
	tPictureAspect      bool
	tBigUserCoords      bool
	tNativeUserCoords   bool
	tFixAspect          bool
	tSequentialNumerics bool
	tRegionNumerics     bool
	tRegionConfig       bool
	tConcurrentConfig   bool
	tDraw               bool

	iterLimit          uint
	divergeLimit       float64
	pictureW           uint
	pictureH           uint
	pictureAspect      float64
	bigUserMin         BigComplex
	bigUserMax         BigComplex
	nativeUserMin      complex128
	nativeUserMax      complex128
	fixAspect          bool
	sequentialNumerics SequentialNumerics
	regionNumerics     RegionNumerics
	regionConfig       RegionParameters
	concurrentConfig   ConcurrentRegionParameters
	draw               DrawingContext
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
