package libgodelbrot

import (
	"image"
)

func Render(info *Info) (*image.NRGBA, error) {
	context, err := MakeRenderer(info)
	if err == nil {
		return context.Render()
	} else {
		return nil, err
	}
}

var __DEBUG = true