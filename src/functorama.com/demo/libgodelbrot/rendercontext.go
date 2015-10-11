package libgodelbrot

import (
    "fmt"
    "log"
    "math/big"
    "math"
    "strconv"
)

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

type ContextMediator struct {
    Info RenderInfo
    Numerics NumericsSystem
    Renderer RenderStrategy
    // The palette
    Palette Palette
}

// Based on the description, choose a renderer, numerical system and palette
// and combine them into a coherent render context
func (desc RenderDescription) CreateInitialRenderContext() RenderContext {
    mediator := &ContextMediator{
        Info: RenderInfo {
            UserDescription: desc
        }
    }

    mediator.initNumerics()
    mediator.initRenderStrategy()
    mediator.initPalette()
}

func (mediator *ContextMediator) initPalette() {
    desc := mediator.Info.UserDescription
    // We are planning more types of palettes soon
    switch desc.PaletteType {
    case StoredPalette:
        mediator.createStoredPalette(desc.PaletteCode)
    default:
        panic(fmt.Sprintf("Unknown palette kind: %v", desc.PaletteType))
    }
}

// Initialize the render system
func (mediator *ContextMediator) initRenderStrategy() {
    desc := mediator.Info.UserDescription
    switch desc.RenderMode {
    case AutoDetectRenderMode:
        mediator.chooseFastRenderStrategy()
    case SequentialRenderMode:
        mediator.useSequentialRenderer()
    case RegionRenderMode:
        mediator.useRegionRenderer()
    case ConcurrentRegionRenderMode:
        mediator.useConcurrentRegionRenderer()
    default:
        panic(fmt.Sprintf("Unknown render mode: %v", desc.RenderMode))
    }
}

// Initialize the numerics system
func (mediator *ContextMediator) initNumerics() {  
    desc := mediator.Info.UserDescription
    mediator.parseUserCoords()
    switch desc.Numerics {
    case AutoDetectNumericsMode:
        mediator.chooseAccurateNumerics()
    case NativeNumericsMode:
        mediator.useNativeNumerics()
    case BigFloatNumericsMode:
        mediator.useBigFloatNumerics()
    default:
        panic(fmt.Sprintf("Unknown numerics mode: %v", desc.Numerics))
    }
}

func (mediator *ContextMediator) createStoredPalette() {
    palettes := map[string]PaletteFactory {
        "redscale": NewRedscalePalette,
        "pretty": NewPrettyPalette,
    }
    code := mediator.Info.UserDescription.PaletteCode
    found := palettes[code]
    if found == nil {
        log.Fatal("Unknown palette: ", code)
    }
    mediator.Palette = found
}

func (mediator *ContextMediator) chooseAccurateNumerics() {
    desc := mediator.Info.UserDescription

    realAccurate := isPixelPerfect(desc.RealMin, desc.RealMax, desc.Width)
    imagAccurate := isPixelPerfect(desc.ImagMin, desc.ImagMax, desc.Height)

    if realAccurate && imagAccurate {
        mediator.useNativeNumerics()
    } else {
        mediator.useBigFloatNumerics()
    }
}

func (mediator *ContextMediator) useNativeNumerics() {
    mediator.Numerics = NewNativeNumerics(mediator)
    mediator.Info.DetectedNumericsMode = NativeNumericsMode
}

func (mediator *ContextMediator) useBigFloatNumerics() {
    mediator.Info.DetectedNumericsMode = BigFloatNumericsMode
    mediator.Numerics = NewBigFloatNumerics(mediator)
}

func (mediator *ContextMediator) parseUserCoords() {
    nativeActions := []func(float64) {
        func (realMin float64) { mediator.Info.NativeRealMin = realMin },
        func (realMax float64) { mediator.Info.NativeRealMax = realMax },
        func (imagMin float64) { mediator.Info.NativeImagMin = imagMin },
        func (imagMax float64) { mediator.Info.NativeImagMax = imagMax },
    }

    bigActions := []func(big.Float) {
        func (realMin big.Float) { mediator.Info.BigRealMin = realMin },
        func (realMax big.Float) { mediator.Info.BigRealMax = realMax },
        func (imagMin big.Float) { mediator.Info.BigImagMin = imagMin },
        func (imagMax big.Float) { mediator.Info.BigImagMax = imagMax },
    }

    desc := mediator.Info.UserDescription
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
            parseFatal(nativeErr ? nativeErr == nil : bigErr, inputNames[i])
        }

        // Handle parse results
        nativeActions[i](native)
        bigActions[i](bigFloat)
    }
}

// Choose an optimal strategy for rendering the image
func (mediator *ContextMediator) chooseFastRenderStrategy() {
    desc := mediator.Info.UserDescription

    area := desc.ImageWidth * desc.ImageHeight
    numerics := mediator.Info.DetectedNumericsMode

    if numerics == AutoDetectNumericsMode {
        panic("Must choose render strategy after numerics system")
    }

    if area < DefaultTinyImageArea && numerics == NativeNumericsMode {
        // Use `SequenceRenderStrategy' when
        // We have native arithmetic and the image is tiny
        mediator.useSequentialRenderer()
    } else if desc.Jobs <= DefaultLowThreading {
        // Use `RegionRenderStrategy' when 
        // the number of jobs is small
        mediator.useRegionRenderer()
    } else { 
        // Use `ConcurrentRegionRenderStrategy' otherwise
        mediator.useConcurrentRegionRenderer()
    }
}

func (mediator *ContextMediator) useSequentialRenderer() {
    mediator.Renderer = NewSequentialRenderer(mediator)
    mediator.Info.DetectedRenderStrategy = SequenceRenderStrategy
}

func (mediator *ContextMediator) useRegionRenderer() {
    mediator.Renderer = NewRegionRenderer(mediator)
    mediator.Info.DetectedRenderStrategy = RegionRenderStrategy
}

func (mediator *ContextMediator) useConcurrentRegionRenderer() {
    mediator.Renderer = NewConcurrentRegionRenderer(mediator)
    mediator.Info.DetectedRenderStrategy
}

// True if we can reperesent the required number of divisions between min and max
func isPixelPerfect(bottom float64, top float64, divisions uint) bool {
    for i := 0; i < int(divisions); i++ {
        bottom = math.Nextafter(bottom, math.MaxFloat64)
    }
    return bottom < top
}

// Trigger a fatal error
func parseFatal(err error, inputName string) {
    log.Fatal("Could not parse '", inputName, "': ", err)
}

// Parse a big.Float
func parseBig(number string) {
    // Do we need to care about the actual base used?
    f, _, err := big.ParseFloat(number, DefaultBase, DefaultHighPrec)
    return f, err
}