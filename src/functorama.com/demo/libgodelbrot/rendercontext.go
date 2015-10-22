package libgodelbrot

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/big"
	"runtime"
	"strconv"
)

// A full Godelbrot render context that can render the fractal to a picture
type RenderContext interface {
	Render() (image.NRGBA, error)
}

// Available kinds of palettes
type PaletteKind uint

const (
	StoredPalette = PaletteKind(iota)
)

// Available render algorithms
type RenderMode uint

const (
	AutoDetectRenderMode       = RenderMode(iota)
	RegionRenderMode           = RenderMode(iota)
	SequentialRenderMode       = RenderMode(iota)
	ConcurrentRegionRenderMode = RenderMode(iota)
)

// Available numeric systems
type NumericsMode uint

const (
	// Functions should auto-detect the correct system for rendering
	AutoDetectNumericsMode = NumericsMode(iota)
	// Use the native CPU arithmetic operations
	NativeNumericsMode = NumericsMode(iota)
	// Use arithmetic based around the standard library big.Float type
	BigFloatNumericsMode = NumericsMode(iota)
)

// User description of the render to be accomplished
type RenderDescription struct {
	IterateLimit uint8
	DivergeLimit float64
	RealMin      string
	RealMax      string
	ImagMin      string
	ImagMax      string
	ImageWidth   uint
	ImageHeight  uint
	// Size of thread input buffer
	ThreadBufferSize uint
	PaletteType      PaletteKind
	PaletteCode      string
	FixAspect        bool
	// Render algorithm
	Renderer RenderMode
	// Number of render threads
	Jobs           uint
	RegionCollapse uint
	// Numerical system
	Numerics NumericsMode
	// Number of samples taken when detecting region render glitches
	GlitchSamples uint
}

func DefaultRenderDescription() *RenderDescription {
	jobs := runtime.NumCPU() - 1
	return &RenderDescription{
		IterateLimit:   DefaultIterations,
		DivergeLimit:   DefaultDivergeLimit,
		RegionCollapse: DefaultCollapse,
		BufferSize:     DefaultBufferSize,
		GlitchSamples:  DefaultGlitchSamples,
		Jobs:           uint(jobs),
		RealMin:        float2str(real(MandelbrotMin)),
		ImagMin:        float2str(imag(MandelbrotMin)),
		RealMax:        float2str(real(MandelbrotMax)),
		ImagMax:        float2str(imag(MandelbrotMax)),
		ImageHeight:    DefaultImageHeight,
		ImageWidth:     DefaultImageWidth,
		PaletteType:    StoredPalette,
		PaletteCode:    "pretty",
	}
}

func float2str(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
