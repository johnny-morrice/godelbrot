package libgodelbrot

// An opaque facade used by subsystems to interact with the application at large 
type RenderApplication interface {
    GlitchSamples() uint
    Limits() (uint, float64)
    RegionCollapseSize() uint
    PictureDimensions() (uint, uint)
    PictureAspect() float64
    BigUserCoords() (BigComplex, BigComplex)
    NativeUserCoords() (complex128, complex128)
    FixAspect() bool
    SequentialNumerics() SequentialNumerics
    RegionNumerics() RegionNumerics
}