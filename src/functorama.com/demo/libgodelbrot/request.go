package libgodelbrot

import (
    "math/big"
    "runtime"
    "strconv"
)

// Info completely describes the render process
type Info struct {
    UserRequest Request
    // Describe the render strategy in use
    RenderStrategy RenderMode
    // Describe the numerics system in use
    NumericsStrategy NumericsMode
    Precision uint
    RealMin big.Float
    RealMax big.Float
    ImagMin big.Float
    ImagMax big.Float
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
    RegionRenderMode
    SequenceRenderMode
    SharedRegionRenderMode
)

// Available numeric systems
type NumericsMode uint

const (
    // Functions should auto-detect the correct system for rendering
    AutoDetectNumericsMode = NumericsMode(iota)
    // Use the native CPU arithmetic operations
    NativeNumericsMode
    // Use arithmetic based around the standard library big.Float type
    BigFloatNumericsMode
)

// Request is a user description of the render to be accomplished
type Request struct {
    IterateLimit uint8
    DivergeLimit float64
    RealMin      string
    RealMax      string
    ImagMin      string
    ImagMax      string
    ImageWidth   uint
    ImageHeight  uint

    PaletteType      PaletteKind
    PaletteCode      string
    FixAspect        bool
    // Render algorithm
    Renderer RenderMode
    // Number of render threads
    Jobs           uint16
    RegionCollapse uint
    // Numerical system
    Numerics NumericsMode
    // Number of samples taken when detecting region render glitches
    GlitchSamples uint
    // Number of bits for big.Float rendering
    Precision uint
}

func DefaultRequest() *Request {
    jobs := runtime.NumCPU() - 1
    return &Request{
        IterateLimit:   DefaultIterations,
        DivergeLimit:   DefaultDivergeLimit,
        RegionCollapse: DefaultCollapse,
        GlitchSamples:  DefaultGlitchSamples,
        Jobs:           uint16(jobs),
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