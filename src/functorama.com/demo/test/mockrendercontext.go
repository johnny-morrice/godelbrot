package test

import (
	"image"
)

type mockRenderContext struct {
	tRender bool
	picture *image.NRGBA
	err     error
}

func (mock *mockRenderContext) Render() (*image.NRGBA, error) {
	mock.tRender = bool
	return mock.picture, mock.err
}
