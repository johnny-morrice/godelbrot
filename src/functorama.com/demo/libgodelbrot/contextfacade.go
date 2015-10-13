package libgodelbrot

// Context facade implements SystemRenderContext
type ContextFacade struct {
    config ContextInit
}

// The context facade implements RenderContext by drawing to an image
func (context *ContextFacade) Render() (image.NRGBA, error) {
    return context.config.RenderContext.Render()
}

func (context *ContextFacade) GlitchSamples() uint {
    return context.config.info.UserDescription.GlitchSamples
}

// Provide the iteration and divergence limits
func (context *ContextFacade) Limits() (uint, float64) {
    desc := context.config.info.UserDescription
    return desc.IterateLimit, desc.DivergeLimit
}

// Provide the region collapse size
func (context *ContextFacade) RegionCollapseSize() uint {
    return context.config.info.UserDescription.RegionCollapse
}

// Provide the image dimensions
func (context *ContextFacade) PictureDimensions() (uint, uint) {
    desc := context.config.info.UserDescription
    return desc.ImageWidth, desc.ImageHeight
}

// Provide the image aspect ratio
func (context *ContextFacade) PictureAspect() float64 {
    pictureWidth, pictureHeight := context.config.PictureDimensions()
    return float64(pictureWidth) / float64(pictureHeight)
}

// Provide the min and max plane coordinates, respectively, as defined by the user
func (context *ContextFacade) BigUserCoords() (BigComplex, BigComplex) {
    info := context.config.info
    min := BigComplex{info.BigRealMin, info.BigImagMin}
    max := BigComplex{info.BigRealMax, info.BigImagMax}
    return min, max
}

// Provide the min and max plane coordinates, respectively, as defined by the user
func (context *ContextFacade) NativeUserCoords() (complex128, complex128) {
    info := context.config.info
    return complex(info.RealMin, info.ImagMin), complex(info.RealMax, info.ImagMax)
}

func (context *ContextFacade) FixAspect() bool {
    return context.config.info.UserDescription.FixAspect
}

func (context *ContextFacade) SequentialNumerics() SequentialNumerics {
    return context.config.numerics
}

func (context *ContextFacade) RegionNumerics() RegionNumerics {
    return context.config.numerics
}