package libgodelbrot

// An interface to the Godelbrot system at large
type GodelbrotApp struct {
    config *ContextInit
}

// Number of samples to be taken when correcting render optimization glitches
func (context *GodelbrotApp) GlitchSamples() uint {
    return context.config.info.UserDescription.GlitchSamples
}

// Provide the iteration and divergence limits
func (context *GodelbrotApp) Limits() (uint, float64) {
    desc := context.config.info.UserDescription
    return desc.IterateLimit, desc.DivergeLimit
}

// Provide the region collapse size
func (context *GodelbrotApp) RegionCollapseSize() uint {
    return context.config.info.UserDescription.RegionCollapse
}

// Provide the image dimensions
func (context *GodelbrotApp) PictureDimensions() (uint, uint) {
    desc := context.config.info.UserDescription
    return desc.ImageWidth, desc.ImageHeight
}

// Provide the image aspect ratio
func (context *GodelbrotApp) PictureAspect() float64 {
    pictureWidth, pictureHeight := context.config.PictureDimensions()
    return float64(pictureWidth) / float64(pictureHeight)
}

// Provide the min and max plane coordinates, respectively, as defined by the user
func (context *GodelbrotApp) BigUserCoords() (BigComplex, BigComplex) {
    info := context.config.info
    min := BigComplex{info.BigRealMin, info.BigImagMin}
    max := BigComplex{info.BigRealMax, info.BigImagMax}
    return min, max
}

// Provide the min and max plane coordinates, respectively, as defined by the user
func (context *GodelbrotApp) NativeUserCoords() (complex128, complex128) {
    info := context.config.info
    return complex(info.RealMin, info.ImagMin), complex(info.RealMax, info.ImagMax)
}

func (context *GodelbrotApp) FixAspect() bool {
    return context.config.info.UserDescription.FixAspect
}

func (context *GodelbrotApp) SequentialNumerics() SequentialNumerics {
    return context.config.numerics
}

func (context *GodelbrotApp) RegionNumerics() RegionNumerics {
    return context.config.numerics
}