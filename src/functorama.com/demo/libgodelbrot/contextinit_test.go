package libgodelbrot

import (
	"fmt"
	"math"
	"testing"
)

func TestNewRenderContextFactory(t *testing.T) {
	if testing.Short {
		t.Skip("Skipping in short mode")
	}

	desc := DefaultRenderDescription()
	goodFactory, noError := NewRenderContextFactory(desc)

	if goodFactory == nil {
		t.Fatal("Good factory was nil")
	}

	if goodFactory.info == nil {
		t.Error("info was nil")
	}

	if noError != nil {
		t.Error("Expected no error")
	}

	// Copy default description
	badDesc := &(*desc)
	badDesc.Renderer = RenderMode(200)
	badFactory, isError := NewRenderContextFactory(badDesc)

	if badFactory != nil {
		t.Error("Bad factory existed")
	}

	if isError == nil {
		t.Error("Expected error")
	}
}

func TestContextInitPalette(t *testing.T) {
	desc := &RenderDescription{
		PaletteType: StoredPalette,
		PaletteCode: "pretty",
	}
	context := &RenderContextFactory{
		info: RenderInfo{
			UserDescription: desc,
		},
	}
	context.initPalette()
	switch context.palette.(type) {
	case PrettyPalette:
		// All good
	default:
		t.Error("Expected type PrettyPalette but received:", context.palette)
	}

	desc.PaletteCode = "redscale"
	context.initPalette()
	switch context.palette.(type) {
	case RedscalePalette:
		// All good
	default:
		t.Error("Expected type RedscalePalette but received:", context.palette)
	}
}

func TestContextInitNumerics(t *testing.T) {
	// First test auto configure
	autoNativeEasy := RenderDescription{
		Numerics:    AutoDetectNumericsMode,
		ImageWidth:  10,
		ImageHeight: 10,
		RealMin:     "0.0",
		RealMax:     "10.0",
		ImagMin:     "0.0",
		ImagMax:     "10.0",
	}

	superLarge := "1"
	zeros := 100
	for i := 0; i < zeros; i++ {
		superLarge = append(superLarge, "0")
	}
	superLarge = append(superLarge, ".0")

	autoBigEasy := autoNativeEasy
	autoBigEasy.ImageWidth = superLarge
	autoBigEasy.ImageHeight = superLarge

	boundaryUnit := math.Nextafter(0.0, 1.0)
	nativeBoundary := float2str(boundaryUnit * 10)
	bigBoundary := float2str(boundaryUnit * 9)

	autoNativeEdge := RenderDescription{
		numerics: AutoDetectNumericsMode,
		RealMin:  "0.0",
		RealMax:  nativeBoundary,
		ImagMin:  "0.0",
		ImagMax:  nativeBoundary,
	}

	autoBigEdge := RenderDescription{
		Numerics: AutoDetectNumericsMode,
		RealMin:  "0.0",
		RealMax:  bigBoundary,
		ImagMin:  "0.0",
		ImagMax:  bigBoundary,
	}

	context := &RenderContextFactory{}

	context.info.UserDescription = autoNativeEasy
	contextNativeCheck(t, context)

	context.info.UserDescription = autoNativeEdge
	contextNativeCheck(t, context)

	context.info.UserDescription = autoBigEasy
	contextBigCheck(t, context)

	context.info.UserDescription = autoBigEdge
	contextBigCheck(t, context)

	// Check manual settings
	context.info.UserDescription = RenderDescription{}

	context.info.UserDescription.Numerics = NativeNumericsMode
	contextNativeCheck(t, context)

	context.info.UserDescription.Numerics = BigNumericsMode
	contextBigCheck(t, context)

}

func contextNativeCheck(t *testing.T, context *RenderContextFactory) {
	context.initNumerics()
	if context.info.DetectedNumericsMode != NativeNumericsMode {
		t.Error("Expected native numerics, but received:",
			context.info.DetectedNumericsMode)
	}
}

func contextBigCheck(t *testing.T, context *RenderContextFactory) {
	context.initNumerics()
	if context.info.DetectedNumericsMode != BigNumericsMode {
		t.Error("Expected big numerics, but received:",
			context.info.DetectedNumericsMode)
	}
}

func TestInitRenderStrategy(t *testing.T) {
	const small = 10
	const large = 400
	const smallJobs = 1
	// olololol
	const bigJobs = 20

	// Check auto first
	autoSequence := &RenderContextFactory{
		info: RenderInfo{
			DetectedNumericsMode: NativeNumericsMode,
			UserDescription: RenderDescription{
				ImageWidth:  small,
				ImageHeight: small,
			},
		},
	}

	autoNativeRegion := &RenderContextFactory{
		info: RenderInfo{
			DetectedNumericsMode: NativeNumericsMode,
			UserDescription: RenderDescription{
				Jobs:        smallJobs,
				ImageWidth:  large,
				ImageHeight: large,
			},
		},
	}

	autoBigRegion := &RenderContextFactory{
		info: RenderInfo{
			DetectedNumericsMode: BigNumericsMode,
			UserDescription: RenderDescription{
				Jobs:        smallJobs,
				ImageWidth:  small,
				ImageHeight: small,
			},
		},
	}

	autoBigConcurrent := &RenderContextFactory{
		info: RenderInfo{
			DetectedNumericsMode: BigNumericsMode,
			UserDescription: RenderDescription{
				Jobs: bigJobs,
			},
		},
	}

	autoNativeConcurrent := &RenderContextFactory{
		info: RenderInfo{
			DetectedNumericsMode: NativeNumericsMode,
			UserDescription: RenderDescription{
				Jobs: bigJobs,
			},
		},
	}

	// Check manual render mode setting
	manual := *&autoNativeRegion

	renderCheck(t, autoSequence, SequenceRenderMode)
	renderCheck(t, autoNativeRegion, RegionRenderMode)
	renderCheck(t, autoBigRegion, RegionRenderMode)
	renderCheck(t, autoBigConcurrent, ConcurrentRenderMode)
	renderCheck(t, autoNativeConcurrent, ConcurrentRenderMode)

	modes := []RenderMode[RegionRenderMode, SequenceRenderMode, ConcurrentRenderMode]
	for _, mode := range modes {
		manual.info.UserDescription.Numerics = mode
		renderCheck(t, manual, mode)
	}
}

func renderCheck(t *testing.T, context *RenderContextFactory, strategy RenderMode) {
	context.initRenderStrategy()
	if context.info.DetectedNumericsMode != strategy {
		t.Error("Expected render strategy", strategy,
			"but received:", context.info.DetectedRenderStrategy)
	}
}
