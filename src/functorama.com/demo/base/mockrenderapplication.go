package base

type MockRenderApplication struct {
	TBaseConfig bool
	TPictureDimensions bool

	IterateLimit          uint8
	DivergeLimit       float64
	FixAspect          bool

	PictureWidth uint
	PictureHeight uint
}

func (mock *MockRenderApplication) BaseConfig() BaseConfig {
	mock.TBaseConfig = true
	return BaseConfig{
		mock.IterateLimit,
		mock.DivergeLimit,
		mock.FixAspect,
	}
}
func (mock *MockRenderApplication) PictureDimensions() (uint, uint) {
	mock.TPictureDimensions = true
	return mock.PictureWidth, mock.PictureHeight
}
