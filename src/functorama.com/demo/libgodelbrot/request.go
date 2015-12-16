package libgodelbrot

import (
    "math/big"
    "runtime"
    "strconv"
)

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
    SequenceRenderMode       = RenderMode(iota)
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

func DefaultRequest() *Request {
    jobs := runtime.NumCPU() - 1
    return &Request{
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