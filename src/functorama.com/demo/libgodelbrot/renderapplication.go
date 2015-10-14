package libgodelbrot

type RegionParameters struct {
    GlitchSamples uint
    RegionCollapseSize uint
}

type ConcurrentRegionParameters struct {
    BufferSize uint
    RenderJobs uint
}

// An opaque facade used by subsystems to interact with the application at large 
type RenderApplication interface {
    Limits() (uint, float64)
    PictureDimensions() (uint, uint)
    PictureAspect() float64
    BigUserCoords() (BigComplex, BigComplex)
    NativeUserCoords() (complex128, complex128)
    FixAspect() bool
    SequentialNumerics() SequentialNumerics
    RegionNumerics() RegionNumerics
    RegionConfig() RegionParameters
    ConcurrentConfig() ConcurrentRegionParameters
    Draw() DrawingContext
}