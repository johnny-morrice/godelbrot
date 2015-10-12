package libgodelbrot

import (
    "fmt"
    "log"
    "math/big"
    "math"
    "strconv"
    "image"
    "runtime"
)

// Something that can render the mandelbrot set
type RenderContext interface {
    Render() (image.Image, error)
}

type PaletteKind uint

const (
    StoredPalette = PaletteKind(iota)
)

type RenderMode uint
const (
    AutoDetectRenderMode = RenderMode(iota)
    RegionRenderMode = RenderMode(iota)
    SequentialRenderMode = RenderMode(iota)
    ConcurrentRegionRenderMode = RenderMode(iota)
)

type NumericsMode uint
const (
    // Functions should auto-detect the correct system for rendering
    AutoDetectNumericsMode = NumericsMode(iota)
    // Use the native CPU arithmetic operations
    NativeNumericsMode = NumericsMode(iota)
    // Use arithmetic based around the standard library big.Float type
    BigFloatNumericsMode = NumericsMode(iota)
)

// High level description of the render to be accomplished
type RenderDescription {
    IterateLimit   uint8
    DivergeLimit   float64
    RealMin string
    RealMax string
    ImagMin string
    ImagMax string
    ImageWidth uint
    ImageHeight uint
    // Size of thread input buffer
    ThreadBufferSize uint
    PaletteType PaletteKind
    PaletteCode string
    FixAspect bool
    // Render algorithm
    Renderer RenderMode
    // Number of render threads
    Jobs uint
    RegionCollapse uint
    BufferSize uint
    // Numerical system
    Numerics NumericsMode
}

// Based on the description, choose a renderer, numerical system and palette
// and combine them into a coherent render context
func (desc RenderDescription) CreateInitialRenderContext() (context RenderContext, err error) {
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

    context = &ContextFacade{
        Info: RenderInfo {
            UserDescription: desc
        }
    }

    context.initNumerics()
    context.initRenderStrategy()
    context.initPalette()

    return
}

type RenderInfo struct {
    UserDescription RenderDescription
    // Describe the render strategy in use
    DetectedRenderStrategy RenderMode
    // Describe the numerics system in use
    DetectedNumericsMode NumericsMode
    // RealMin as a native float
    NativeRealMin float64
    // RealMax as a native float
    NativeRealMax float64
    // ImagMin as a native float
    NativeImagMin float64
    // ImagMax as a native float
    NativeImagMax float64
    // RealMin as a big float
    BigRealMin big.Float
    // RealMax as a big float
    BigRealMax big.Float
    // ImagMin as a big float
    BigImagMin big.Float
    // ImagMax as a big float
    BigImagMax big.Float
}

type ContextFacade struct {
    Info RenderInfo
    Numerics NumericsSystem
    Renderer RenderContext
    // The palette
    Palette Palette
}

// Provide the iteration and divergence limits
func (context *ContextFacade) Limits() (uint, float64) {
    desc := context.Info.UserDescription
    return desc.IterateLimit, desc.DivergeLimit
}

// Provide the region collapse size
func (context *ContextFacade) RegionCollapseSize() uint {
    return context.Info.UserDescription.RegionCollapse
}

// Provide the image dimensions
func (context *ContextFacade) PictureDimensions() (uint, uint) {
    desc := context.Info.UserDescription
    return desc.ImageWidth, desc.ImageHeight
}

// Provide the min and max plane coordinates, respectively, as defined by the user
func (context *ContextFacade) NativeUserCoords() (complex128, complex128) {
    info := context.Info
    return complex(info.RealMin, info.RealMax), complex(info.ImagMin, info.ImagMax)
}

func (context *ContextFacade) FixAspect() bool {
    return context.Info.UserDescription.FixAspect
}

func (context *ContextFacade) SequentialNumerics() SequentialNumerics {
    return context.Numerics
}

func (context *ContextFacade) RegionNumerics() RegionNumerics {
    return context.Numerics
}

func (context *ContextFacade) initPalette() {
    desc := context.Info.UserDescription
    // We are planning more types of palettes soon
    switch desc.PaletteType {
    case StoredPalette:
        context.createStoredPalette(desc.PaletteCode)
    default:
        panic(fmt.Sprintf("Unknown palette kind: %v", desc.PaletteType))
    }
}

// Initialize the render system
func (context *ContextFacade) initRenderStrategy() {
    desc := context.Info.UserDescription
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
func (context *ContextFacade) initNumerics() {  
    desc := context.Info.UserDescription
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

func (context *ContextFacade) createStoredPalette() {
    palettes := map[string]PaletteFactory {
        "redscale": NewRedscalePalette,
        "pretty": NewPrettyPalette,
    }
    code := context.Info.UserDescription.PaletteCode
    found := palettes[code]
    if found == nil {
        log.Fatal("Unknown palette: ", code)
    }
    context.Palette = found
}

func (context *ContextFacade) chooseAccurateNumerics() {
    desc := context.Info.UserDescription

    realAccurate := isPixelPerfect(desc.RealMin, desc.RealMax, desc.Width)
    imagAccurate := isPixelPerfect(desc.ImagMin, desc.ImagMax, desc.Height)

    if realAccurate && imagAccurate {
        context.useNativeNumerics()
    } else {
        context.useBigFloatNumerics()
    }
}

func (context *ContextFacade) useNativeNumerics() {
    context.Numerics = NewNativeNumerics(context)
    context.Info.DetectedNumericsMode = NativeNumericsMode
}

func (context *ContextFacade) useBigFloatNumerics() {
    context.Info.DetectedNumericsMode = BigFloatNumericsMode
    context.Numerics = NewBigFloatNumerics(context)
}

func (context *ContextFacade) parseUserCoords() {
    nativeActions := []func(float64) {
        func (realMin float64) { context.Info.NativeRealMin = realMin },
        func (realMax float64) { context.Info.NativeRealMax = realMax },
        func (imagMin float64) { context.Info.NativeImagMin = imagMin },
        func (imagMax float64) { context.Info.NativeImagMax = imagMax },
    }

    bigActions := []func(big.Float) {
        func (realMin big.Float) { context.Info.BigRealMin = realMin },
        func (realMax big.Float) { context.Info.BigRealMax = realMax },
        func (imagMin big.Float) { context.Info.BigImagMin = imagMin },
        func (imagMax big.Float) { context.Info.BigImagMax = imagMax },
    }

    desc := context.Info.UserDescription
    userInput := []string {
        desc.RealMin,
        desc.RealMax,
        desc.ImagMin,
        desc.ImagMax,
    }

    inputNames := []string {"realMin", "realMax", "imagMin", "imagMax"}

    for i, num := range(userInput) {
        // Parse a float64 from `num' into `native'
        bits := 64
        native, nativeErr := strconv.ParseFloat(num, bits)
        bigFloat, bigErr := parseBig(num)

        // Handle errors by vomiting organs
        if nativeErr != nil || bigErr != nil {
            parsePanic(nativeErr ? nativeErr == nil : bigErr, inputNames[i])
        }

        // Handle parse results
        nativeActions[i](native)
        bigActions[i](bigFloat)
    }
}

// Choose an optimal strategy for rendering the image
func (context *ContextFacade) chooseFastRenderStrategy() {
    desc := context.Info.UserDescription

    area := desc.ImageWidth * desc.ImageHeight
    numerics := context.Info.DetectedNumericsMode

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

func (context *ContextFacade) useSequentialRenderer() {
    context.Renderer = NewSequentialRenderer(context)
    context.Info.DetectedRenderStrategy = SequenceRenderMode
}

func (context *ContextFacade) useRegionRenderer() {
    context.Renderer = NewRegionRenderer(context)
    context.Info.DetectedRenderStrategy = RegionRenderMode
}

func (context *ContextFacade) useConcurrentRegionRenderer() {
    context.Renderer = NewConcurrentRegionRenderer(context)
    context.Info.DetectedRenderStrategy = ConcurrentRegionRenderMode
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