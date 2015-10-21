package libgodelbrot

import (
    "image"
)

type RegionParameters struct {
    GlitchSamples uint
    RegionCollapseSize uint
}

type ConcurrentRegionParameters struct {
    BufferSize uint
    RenderJobs uint
}

// A facade used by subsystems to interact with the application at large 
type RenderApplication interface {
    // Basic configuration
    IterateLimit() uint8
    DivergeLimit() float64
    PictureDimensions() (uint, uint)
    BigUserCoords() (BigComplex, BigComplex)
    NativeUserCoords() (complex128, complex128)
    FixAspect() bool

    // Configuration for particular render strategies
    RegionConfig() RegionParameters
    ConcurrentConfig() ConcurrentRegionParameters

    // Views into the numerics system used by various render strategies
    RegionNumerics() RegionNumerics
    SequentialNumerics() SequentialNumerics
    SharedRegionNumerics() SharedRegionNumerics
    SharedSequentialNumerics() SharedSequentialNumerics

    // Image drawing facilities
    Draw() DrawingContext
}