package libgodelbrot

import (
	"image"
)

// A facade for users to interact with the configured system
type GodelbrotUserFacade struct {
	Info RenderInfo
	context RenderContext
}

// The user facade will render an image
func (facade *GodelbrotUserFacade) Render() (*image.NRGBA, error) {
	return context.Render()
}

// Create a simple facade for clients to interface with the Godelbrot system
func NewUserFacade(factory *RenderContextFactory) *GodelbrotUserFacade {
	return &GodelbrotUserFacade{
		context: factory.Build(),
		Info: factory.Info,
	}
}

// Object to initialize the godelbrot system
type RenderContextFactory struct {
	info     RenderInfo
}

// InitializeContext examines the description, chooses a renderer, numerical system and palette.
// Together these form a coherent render factory.
func NewRenderContextFactory(desc *RenderDescription) (*RenderContextFactory, error) {
	anything, err := panic2err(func() interface{} {
		factory := &RenderContextFactory{}
		factory.info.UserDescription = *desc

		factory.chooseNumerics()
		factory.chooseRenderStrategy()

		return factory
	})

	factory, ok := anything.(*RenderContext)
	// Explicitly note this is a bug when it fails
	if !ok {
		log.Fatal("BUG on type conversion:", anything)
	}

	return factory
}

func (factory *RenderContextFactory) Build() (RenderContext, error) {
	anything, err := panic2err(func() interface{} {
		switch factory.info.DetectedRenderStrategy {
		case SequentialRenderMode:
			return sequence.NewSequentialRenderer(NewSequenceFacade(factory.info))
		case RegionRenderMode:
			return region.NewRegionRenderer(NewRegionFacade(factory.info))
		case ConcurrentRegionRenderMode:
			return sharedregion.NewSharedRegionRenderer(NewSharedRegionFacade(factory.info))
		default:
			// panic2err?  case in point
			log.Panic("Unsupported render mode:", factory.info.DetectedRenderStrategy)
		}
	})

	context, ok := anything.(RenderContext)
	if !ok {
		log.Fatal("BUG in type conversion:", anything)
	}

	return context, err
}

// Initialize the render system
func (factory *RenderContextFactory) chooseRenderStrategy() {
	desc := factory.info.UserDescription
	switch desc.RenderMode {
	case AutoDetectRenderMode:
		factory.chooseFastRenderStrategy()
	case SequentialRenderMode:
		factory.useSequentialRenderer()
	case RegionRenderMode:
		factory.useRegionRenderer()
	case ConcurrentRegionRenderMode:
		factory.useConcurrentRegionRenderer()
	default:
		log.Panic("Unknown render mode:", desc.RenderMode)
	}
}

// Initialize the numerics system
func (factory *RenderContextFactory) chooseNumerics() {
	desc := factory.info.UserDescription
	factory.parseUserCoords()
	switch desc.Numerics {
	case AutoDetectNumericsMode:
		factory.chooseAccurateNumerics()
	case NativeNumericsMode:
		factory.useNativeNumerics()
	case BigFloatNumericsMode:
		factory.useBigFloatNumerics()
	default:
		log.Panic("Unknown numerics mode:", desc.Numerics)
	}
}


func (factory *RenderContextFactory) chooseAccurateNumerics() {
	desc := factory.info.UserDescription

	realAccurate := isPixelPerfect(desc.RealMin, desc.RealMax, desc.Width)
	imagAccurate := isPixelPerfect(desc.ImagMin, desc.ImagMax, desc.Height)

	if realAccurate && imagAccurate {
		factory.useNativeNumerics()
	} else {
		factory.useBigFloatNumerics()
	}
}

func (factory *RenderContextFactory) useNativeNumerics() {
	factory.numerics = NewNativeNumerics(factory)
	factory.info.DetectedNumericsMode = NativeNumericsMode
}

func (factory *RenderContextFactory) useBigFloatNumerics() {
	factory.info.DetectedNumericsMode = BigFloatNumericsMode
	factory.numerics = NewBigFloatNumerics(factory)
}

func (factory *RenderContextFactory) parseUserCoords() {
	nativeActions := []func(float64){
		func(realMin float64) { factory.info.NativeRealMin = realMin },
		func(realMax float64) { factory.info.NativeRealMax = realMax },
		func(imagMin float64) { factory.info.NativeImagMin = imagMin },
		func(imagMax float64) { factory.info.NativeImagMax = imagMax },
	}

	bigActions := []func(big.Float){
		func(realMin big.Float) { factory.info.BigRealMin = realMin },
		func(realMax big.Float) { factory.info.BigRealMax = realMax },
		func(imagMin big.Float) { factory.info.BigImagMin = imagMin },
		func(imagMax big.Float) { factory.info.BigImagMax = imagMax },
	}

	desc := factory.info.UserDescription
	userInput := []string{
		desc.RealMin,
		desc.RealMax,
		desc.ImagMin,
		desc.ImagMax,
	}

	inputNames := []string{"realMin", "realMax", "imagMin", "imagMax"}

	for i, num := range userInput {
		// Parse a float64 from `num' into `native'
		const bits uint = 64
		native, nativeErr := strconv.ParseFloat(num, bits)
		bigFloat, bigErr := parseBig(num)

		badName := inputNames[i]
		// Handle errors by vomiting organs
		if nativeErr != nil {
			parsePanic(nativeErr, badName)
		}

		if bigErr != nil {
			parsePanic(bigErr, badName)
		}

		// Handle parse results
		nativeActions[i](native)
		bigActions[i](bigFloat)
	}
}

// Choose an optimal strategy for rendering the image
func (factory *RenderContextFactory) chooseFastRenderStrategy() {
	desc := factory.info.UserDescription

	area := desc.ImageWidth * desc.ImageHeight
	numerics := factory.info.DetectedNumericsMode

	if numerics == AutoDetectNumericsMode {
		log.Panic("Must choose render strategy after numerics system")
	}

	if area < DefaultTinyImageArea && numerics == NativeNumericsMode {
		// Use `SequenceRenderStrategy' when
		// We have native arithmetic and the image is tiny
		factory.useSequentialRenderer()
	} else if desc.Jobs <= DefaultLowThreading {
		// Use `RegionRenderStrategy' when
		// the number of jobs is small
		factory.useRegionRenderer()
	} else {
		// Use `ConcurrentRegionRenderStrategy' otherwise
		factory.useConcurrentRegionRenderer()
	}
}

func (factory *RenderContextFactory) useSequentialRenderer() {
	factory.info.DetectedRenderStrategy = SequenceRenderMode
}

func (factory *RenderContextFactory) useRegionRenderer() {
	factory.info.DetectedRenderStrategy = RegionRenderMode
}

func (factory *RenderContextFactory) useConcurrentRegionRenderer() {
	factory.info.DetectedRenderStrategy = ConcurrentRegionRenderMode
}

// True if we can reperesent the required number of divisions between min and max
func isPixelPerfect(bottom float64, top float64, divisions uint) bool {
	for i := 0; i < int(divisions); i++ {
		bottom = math.Nextafter(bottom, math.MaxFloat64)
	}
	return bottom < top
}
