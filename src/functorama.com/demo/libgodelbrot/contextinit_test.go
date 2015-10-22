package libgodelbrot

import (
	"fmt"
	"math"
	"testing"
)

func TestInitializeContext(t *testing.T) {
	desc := DefaultRenderDescription()
	context := InitializeContext(desc)

	if context.info == nil {
		t.Error("info was nil")
	}

	if context.numerics == nil {
		t.Error("info was nil")
	}

	if context.renderer == nil {
		t.Error("renderer was nil")
	}

	if context.palette == nil {
		t.Error("palette was nil")
	}

	if context.picture == nil {
		t.Error("picture was nil")
	}
}

func TestUserFacade(t *testing.T) {
	mock := &mockRenderContext{
		picture: image.NewNRGBA(image.ZR),
		err:     fmt.Errorf("Fake error!"),
	}
	context := ContextInit{RenderContext: mock}
	facade := context.NewUserFacade()
	pic, err := facade.Render()

	if !mock.tRender {
		t.Error("Expected method not called on mock")
	}

	if pic != mock.picture {
		t.Error("Facade returned unexpected picture:", pic)
	}

	if err != mock.err {
		t.Error("Facade returned unexpected error:", err)
	}
}

func TestContextInitPalette(t *testing.T) {
	desc := &RenderDescription{
		PaletteType: StoredPalette,
		PaletteCode: "pretty",
	}
	context := ContextInit{
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

	context := ContextInit{}

	context.info.UserDescription = autoNativeEasy
	contextNativeCheck(t, context, "autoNativeEasy")

	context.info.UserDescription = autoNativeEdge
	contextNativeCheck(t, context, "autoNativeEdge")

	context.info.UserDescription = autoBigEasy
	contextBigCheck(t, context, "autoBigEasy")

	context.info.UserDescription = autoBigEdge
	contextBigCheck(t, context, "autoBigEdge")

	// Check manual settings
	context.info.UserDescription = RenderDescription{}

	context.info.UserDescription.Numerics.NativeNumericsMode
	contextNativeCheck(t, context, "Manual Numerics")

	context.info.UserDescription.Numerics.BigNumericsMode
	contextBigCheck(t, context, "Big Numerics")

}

func contextNativeCheck(t *testing.T, context ContextInit, description string) {
	context.initNumerics()
	if context.info.DetectedNumericsMode != NativeNumericsMode {
		t.Error("When checking ", description,
			"expected metadata to indicate native numerics, but received:",
			context.info.DetectedNumericsMode)
	}

	switch context.numerics.(type) {
	case NativeNumericsFactory:
		// All good
	default:
		t.Error("When checking", description,
			"Expected numerics to be native but received:", context.numerics)
	}
}

func contextBigCheck(t *testing.T, context ContextInit, description string) {
	context.initNumerics()
	if context.info.DetectedNumericsMode != BigNumericsMode {
		t.Error("When checking ", description,
			"expected metadata to indicate native numerics, but received:",
			context.info.DetectedNumericsMode)
	}

	switch context.numerics.(type) {
	case BigNumericsFactory:
		// All good
	default:
		t.Error("When checking", description,
			"Expected numerics to be native but received:", context.numerics)
	}
}

func TestInitRenderStrategy(t *testing.T) {
	const small = 10
	const large = 400
	const smallJobs = 1
	// olololol
	const bigJobs = 20

	// Check auto first
	autoSequence := ContextInit{
		info: RenderInfo{
			DetectedNumericsMode: NativeNumericsMode,
			UserDescription: RenderDescription{
				ImageWidth:  small,
				ImageHeight: small,
			},
		},
	}

	autoNativeRegion := ContextInit{
		info: RenderInfo{
			DetectedNumericsMode: NativeNumericsMode,
			UserDescription: RenderDescription{
				Jobs:        smallJobs,
				ImageWidth:  large,
				ImageHeight: large,
			},
		},
	}

	autoBigRegion := ContextInit{
		info: RenderInfo{
			DetectedNumericsMode: BigNumericsMode,
			UserDescription: RenderDescription{
				Jobs:        smallJobs,
				ImageWidth:  small,
				ImageHeight: small,
			},
		},
	}

	autoBigConcurrent := ContextInit{
		info: RenderInfo{
			DetectedNumericsMode: BigNumericsMode,
			UserDescription: RenderDescription{
				Jobs: bigJobs,
			},
		},
	}

	autoNativeConcurrent := ContextInit{
		info: RenderInfo{
			DetectedNumericsMode: NativeNumericsMode,
			UserDescription: RenderDescription{
				Jobs: bigJobs,
			},
		},
	}
}
