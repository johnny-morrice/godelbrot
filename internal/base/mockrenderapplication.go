package base

type MockRenderApplication struct {
	TBaseConfig        bool
	TPictureDimensions bool

	Base BaseConfig

	PictureWidth  uint
	PictureHeight uint
}

func (mock *MockRenderApplication) BaseConfig() BaseConfig {
	mock.TBaseConfig = true
	return mock.Base
}
func (mock *MockRenderApplication) PictureDimensions() (uint, uint) {
	mock.TPictureDimensions = true
	return mock.PictureWidth, mock.PictureHeight
}
