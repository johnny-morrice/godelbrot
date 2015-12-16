package libgodelbrot

import (
	"image"
)

// Draw the Mandelbrot set.  This is the main entry point to libgodelbrot
func AutoRender(desc *Request) (*image.NRGBA, error) {
	info, configErr := Autoconf(desc)

	if configErr == nil {
		return Render(info)
	} else {
		return nil, configErr
	}
}

func Render(info Info) (*image.NRGBA, error) {
	context, err := Renderer(info)

	if err == nil {
		return context.Render()
	} else {
		return nil, err
	}
}

func AutoConf(desc *Request) (*Info, error) {
	// Autoconf is a thin wrapper, we just pass on to the library internals
	return configure(desc)
}


func AutoZoomNext(fractal *image.NRGBA, previous *Info, rate float64) (*Info, error) {
	bounds, err := zoomDetail(fractal, rate)

	if err == nil {
		return zoomIn(previous, bounds), nil
	} else {
		return nil, err
	}
}

func Renderer(desc *Info) (Renderer,  error) {
	// Renderer is a thin wrapper, we just pass on to the library internals
	return context(desc)
}

type MovieFrame struct {
	ok bool
	still *image.NRGBA
	info *Info
	err   error
}

type autoZoomMovie struct {
	rate float64
	start *Info
	frames chan int
	reel chan MovieFrame
}

type RecordReel interface {
	Wind(frameCount int)
	Iter() <-chan MovieFrame
}

func AutoZoomMovie(rate float64, info *Info) RecordReel {
	reel := autoZoomMovie{}
	reel.rate = rate
	reel.start = info
	reel.frames = make(chan int)
	reel.reel = make(chan MovieFrame)
	return reel
}

func (azm autoZoomMovie) Wind(frameCount int) {
	go func () {
		azm.frames<- frameCount
	}()
}

func (azm autoZoomMovie) Iter() <-chan MovieFrame {
	go func() {
		framesRemaining := <-azm.frames
		info := azm.start
		for i := framesRemaining; i > 0; i-- {
			select {
				// Allow user to request more frames while iterator is running
			case moreFrames := <-azm.frames:
				i += moreFrames
			default:
				// Nothing to do if no extra frames given
			}

			fractal, err := Render(info)
			frame := MovieFrame{
				ok: err != nil,
				still: fractal,
				info: info,
				err: err,
			}
			azm.reel<- frame

			var zoomErr error
			info, zoomErr = AutoZoomNext(fractal, info, azm.rate)

			// If there is an error on zooming, shutdown the iterator
			if zoomErr != nil {
				finalFrame := MovieFrame{
					ok: false,
					still: nil,
					err: zoomErr,
				}
				azm.reel<- finalFrame
				close(azm.reel)
			}
		}
	}()
	return azm.reel
}