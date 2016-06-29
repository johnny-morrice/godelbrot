package godelbrot

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

// These flags are intended for developers only
var __DEBUG = false
var __TRACE = false