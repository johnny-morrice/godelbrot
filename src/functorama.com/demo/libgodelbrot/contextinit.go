package libgodelbrot

import (
	"image"
)

// Object to initialize the godelbrot system
type ContextInit struct {
	info     RenderInfo
	numerics NumericSystem
	renderer RenderContext
	// The palette
	palette Palette
	// The image upon which to draw image
	picture *image.NRGBA
}

// Based on the description, choose a renderer, numerical system and palette
// and combine them into a coherent render context
func InitializeContext(desc RenderDescription) (context ContextInit, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case runtime.Error:
				panic(r)
			default:
				err = r.(error)
			}
		}
	}()

	context = &ContextInit{
		info: RenderInfo{
			UserDescription: desc,
		},
	}

	context.initNumerics()
	context.initRenderStrategy()
	context.initPalette()
	context.initImage()

	return
}

// A facade for users to interact with the configured system
type GodelbrotUserFacade struct {
	config *ContextInit
}

// The user facade will render an image
func (facade *GodelbrotUserFacade) Render() (image.NRGBA, error) {
	return facade.config.RenderContext.Render()
}

// Create a simple facade for clients to interface with the Godelbrot system
func (context *ContextInit) NewUserFacade() *GodelbrotUserFacade {
	return &GodelbrotUserFacade{config: context}
}

// Create a simple facade for subsystems to interface with larger Godelbrot
// system
func (context *ContextInit) NewInnerFacade() GodelbrotApp {
	return GodelbrotApp(context)
}

func (context *ContextInit) initPalette() {
	desc := context.info.UserDescription
	// We are planning more types of palettes soon
	switch desc.PaletteType {
	case StoredPalette:
		context.createStoredPalette(desc.PaletteCode)
	default:
		panic(fmt.Sprintf("Unknown palette kind: %v", desc.PaletteType))
	}
}

// Initialize the render system
func (context *ContextInit) initRenderStrategy() {
	desc := context.info.UserDescription
	switch desc.RenderMode {
	case AutoDetectRenderMode:
		context.chooseFastRenderStrategy()
	case SequentialRenderMode:
		context.useSequentialRenderer()
	case RegionRenderMode:
		context.useRegionRenderer()
	case ConcurrentRegionRenderMode:
		context.useConcurrentRegionRenderer()
	default:
		panic(fmt.Sprintf("Unknown render mode: %v", desc.RenderMode))
	}
}

// Initialize the numerics system
func (context *ContextInit) initNumerics() {
	desc := context.info.UserDescription
	context.parseUserCoords()
	switch desc.Numerics {
	case AutoDetectNumericsMode:
		context.chooseAccurateNumerics()
	case NativeNumericsMode:
		context.useNativeNumerics()
	case BigFloatNumericsMode:
		context.useBigFloatNumerics()
	default:
		panic(fmt.Sprintf("Unknown numerics mode: %v", desc.Numerics))
	}
}

func (context *ContextInit) createStoredPalette() {
	palettes := map[string]PaletteFactory{
		"redscale": NewRedscalePalette,
		"pretty":   NewPrettyPalette,
	}
	code := context.info.UserDescription.PaletteCode
	found := palettes[code]
	if found == nil {
		log.Fatal("Unknown palette: ", code)
	}
	context.palette = found
}

func (context *ContextInit) chooseAccurateNumerics() {
	desc := context.info.UserDescription

	realAccurate := isPixelPerfect(desc.RealMin, desc.RealMax, desc.Width)
	imagAccurate := isPixelPerfect(desc.ImagMin, desc.ImagMax, desc.Height)

	if realAccurate && imagAccurate {
		context.useNativeNumerics()
	} else {
		context.useBigFloatNumerics()
	}
}

func (context *ContextInit) useNativeNumerics() {
	context.numerics = NewNativeNumerics(context)
	context.info.DetectedNumericsMode = NativeNumericsMode
}

func (context *ContextInit) useBigFloatNumerics() {
	context.info.DetectedNumericsMode = BigFloatNumericsMode
	context.numerics = NewBigFloatNumerics(context)
}

func (context *ContextInit) parseUserCoords() {
	nativeActions := []func(float64){
		func(realMin float64) { context.info.NativeRealMin = realMin },
		func(realMax float64) { context.info.NativeRealMax = realMax },
		func(imagMin float64) { context.info.NativeImagMin = imagMin },
		func(imagMax float64) { context.info.NativeImagMax = imagMax },
	}

	bigActions := []func(big.Float){
		func(realMin big.Float) { context.info.BigRealMin = realMin },
		func(realMax big.Float) { context.info.BigRealMax = realMax },
		func(imagMin big.Float) { context.info.BigImagMin = imagMin },
		func(imagMax big.Float) { context.info.BigImagMax = imagMax },
	}

	desc := context.info.UserDescription
	userInput := []string{
		desc.RealMin,
		desc.RealMax,
		desc.ImagMin,
		desc.ImagMax,
	}

	inputNames := []string{"realMin", "realMax", "imagMin", "imagMax"}

	for i, num := range userInput {
		// Parse a float64 from `num' into `native'
		bits := 64
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
func (context *ContextInit) chooseFastRenderStrategy() {
	desc := context.info.UserDescription

	area := desc.ImageWidth * desc.ImageHeight
	numerics := context.info.DetectedNumericsMode

	if numerics == AutoDetectNumericsMode {
		panic("Must choose render strategy after numerics system")
	}

	if area < DefaultTinyImageArea && numerics == NativeNumericsMode {
		// Use `SequenceRenderStrategy' when
		// We have native arithmetic and the image is tiny
		context.useSequentialRenderer()
	} else if desc.Jobs <= DefaultLowThreading {
		// Use `RegionRenderStrategy' when
		// the number of jobs is small
		context.useRegionRenderer()
	} else {
		// Use `ConcurrentRegionRenderStrategy' otherwise
		context.useConcurrentRegionRenderer()
	}
}

func (context *ContextInit) useSequentialRenderer() {
	context.renderer = NewSequentialRenderer(context.NewInnerFacade())
	context.info.DetectedRenderStrategy = SequenceRenderMode
}

func (context *ContextInit) useRegionRenderer() {
	context.renderer = NewRegionRenderer(context.NewInnerFacade())
	context.info.DetectedRenderStrategy = RegionRenderMode
}

func (context *ContextInit) useConcurrentRegionRenderer() {
	context.renderer = NewConcurrentRegionRenderer(context.NewInnerFacade())
	context.info.DetectedRenderStrategy = ConcurrentRegionRenderMode
}

func (context *ContextInit) initPicture() {
	desc := context.info.Desc
	bounds := image.Rectangle{
		Min: image.ZP,
		Max: image.Point{
			X: desc.ImageWidth,
			Y: desc.ImageHeight,
		},
	}
	context.picture = image.NewNRGBA(bounds)
}

// True if we can reperesent the required number of divisions between min and max
func isPixelPerfect(bottom float64, top float64, divisions uint) bool {
	for i := 0; i < int(divisions); i++ {
		bottom = math.Nextafter(bottom, math.MaxFloat64)
	}
	return bottom < top
}

// Panic to escape parsing
func parsePanic(err error, inputName string) {
	return panic(fmt.Sprintf("Could not parse %v: %v'", inputName, err))
}

// Parse a big.Float
func parseBig(number string) {
	// Do we need to care about the actual base used?
	f, _, err := big.ParseFloat(number, DefaultBase, DefaultHighPrec)
	return f, err
}
