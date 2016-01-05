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

func AutoConf(req *Request) (*Info, error) {
	// Configure uses panic when it encounters an error condition.
	// Here we detect that panic and convert it to an error,
	// which is idiomatic for the API.
	anything, err := panic2err(func() interface{} {
		return configure(req)
	})

    if err == nil {
        return anything.(*Info), nil
    } else {
        return nil, err
    }
}

func MakeRenderer(desc *Info) (Renderer,  error) {
	// Renderer is a thin wrapper, we just pass on to the library internals
	return renderer(desc)
}
