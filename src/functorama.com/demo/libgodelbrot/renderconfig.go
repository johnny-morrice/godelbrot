package libgodelbrot

import (
	"fmt"
	"image"
	"math"
	"runtime"
)

type CoordFrame uint

// Co-ordinate frames
const (
	CornerFrame = CoordFrame(iota)
	ZoomFrame   = CoordFrame(iota)
)

// User input
type RenderParameters struct {
	IterateLimit   uint8
	DivergeLimit   float64
	Width          uint
	Height         uint
	Zoom           float64
	RegionCollapse uint
	// Co-ordinate frames
	Frame CoordFrame
	// Top left of view onto plane
	TopLeft complex128
	// Optional Bottom right corner
	BottomRight complex128
	// Number of render threads
	RenderThreads uint
	// Size of thread input buffer
	BufferSize uint
	// Fix aspect
	FixAspect bool
}

func (config RenderParameters) PlaneTopLeft() complex128 {
	return config.TopLeft
}

// Top right of window onto complex plane
func (config RenderParameters) PlaneBottomRight() complex128 {
	if config.Frame == ZoomFrame {
		scaled := MagicSetSize * complex(config.Zoom, 0)
		topLeft := config.PlaneTopLeft()
		right := real(topLeft) + real(scaled)
		bottom := imag(topLeft) - imag(scaled)
		return complex(right, bottom)
	} else if config.Frame == CornerFrame {
		return config.BottomRight
	} else {
		config.framePanic()
	}
	panic("Bug")
	return 0
}

func (config RenderParameters) framePanic() {
	panic(fmt.Sprintf("Unknown frame: %v", config.Frame))
}

func (config RenderParameters) BlankImage() *image.NRGBA {
	return image.NewNRGBA(image.Rectangle{
		Min: image.ZP,
		Max: image.Point{
			X: int(config.Width),
			Y: int(config.Height),
		},
	})
}
func (args RenderParameters) PlaneSize() complex128 {
	if args.Frame == ZoomFrame {
		return complex(args.Zoom, 0) * MagicSetSize
	} else if args.Frame == CornerFrame {
		tl := args.TopLeft
		br := args.BottomRight
		return complex(real(br)-real(tl), imag(tl)-imag(br))
	} else {
		args.framePanic()
	}
	panic("Bug")
	return 0
}

// Configure the render parameters into something working
// Fixes aspect ratio
func (args RenderParameters) Configure() *RenderConfig {
	planeSize := args.PlaneSize()
	planeWidth := real(planeSize)
	planeHeight := imag(planeSize)

	imageAspect := float64(args.Width) / float64(args.Height)
	planeAspect := planeWidth / planeHeight

	if args.FixAspect {
		tl := args.PlaneTopLeft()
		// If the plane aspect is greater than image aspect
		// Then the plane is too short, so must be made taller
		if planeAspect > imageAspect {
			taller := planeWidth / imageAspect
			br := tl + complex(planeWidth, -taller)
			args.BottomRight = br
			args.Frame = CornerFrame
		} else if planeAspect < imageAspect {
			// If the plane aspect is less than the image aspect
			// Then the plane is too thin, and must be made fatter
			fatter := planeHeight * imageAspect
			br := tl + complex(fatter, -planeHeight)
			args.BottomRight = br
			args.Frame = CornerFrame
		}
	}

	return &RenderConfig{
		RenderParameters: args,
		HorizUnit:        planeWidth / float64(args.Width),
		VerticalUnit:     planeHeight / float64(args.Height),
		ImageLeft:        0,
		ImageTop:         0,
	}
}

// Machine prepared input, caching interim results
type RenderConfig struct {
	RenderParameters
	// One pixel's space on the plane
	HorizUnit    float64
	VerticalUnit float64
	ImageLeft    uint
	ImageTop     uint
}

func DefaultRenderThreads() uint {
	cpus := runtime.NumCPU()
	var threads uint
	if cpus > 1 {
		threads = uint(cpus - 1)
	} else {
		threads = 1
	}
	return threads
}

// Use magic values to create default config
func DefaultConfig() *RenderConfig {
	threads := DefaultRenderThreads()
	params := RenderParameters{
		IterateLimit:   DefaultIterations,
		DivergeLimit:   DefaultDivergeLimit,
		Width:          DefaultImageWidth,
		Height:         DefaultImageHeight,
		TopLeft:        MagicOffset,
		Zoom:           DefaultZoom,
		Frame:          ZoomFrame,
		RegionCollapse: DefaultCollapse,
		RenderThreads:  threads,
		BufferSize:     DefaultBufferSize,
	}
	return params.Configure()
}

func (config RenderConfig) PlaneToPixel(c complex128) (rx uint, ry uint) {
	// Translate x
	tx := real(c) - real(config.TopLeft)
	// Scale x
	sx := tx / config.HorizUnit

	// Translate y
	ty := imag(c) - imag(config.TopLeft)
	// Scale y
	sy := ty / config.VerticalUnit

	rx = uint(math.Floor(sx))
	// Remember that we draw downwards
	ry = uint(math.Ceil(-sy))

	return
}

func (config RenderConfig) ImageTopLeft() (uint, uint) {
	return config.ImageLeft, config.ImageTop
}
