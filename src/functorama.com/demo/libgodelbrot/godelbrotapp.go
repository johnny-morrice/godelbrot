package libgodelbrot

// An interface to the Godelbrot system at large
type GodelbrotApp *ContextInit

// Provide the iteration and divergence limits
func (app *GodelbrotApp) Limits() (uint, float64) {
    desc := app.info.UserDescription
    return desc.IterateLimit, desc.DivergeLimit
}

// Provide the region collapse size
func (app *GodelbrotApp) RegionCollapseSize() uint {
    return app.info.UserDescription.RegionCollapse
}

// Provide the image dimensions
func (app *GodelbrotApp) PictureDimensions() (uint, uint) {
    desc := app.info.UserDescription
    return desc.ImageWidth, desc.ImageHeight
}

// Provide the image aspect ratio
func (app *GodelbrotApp) PictureAspect() float64 {
    pictureWidth, pictureHeight := app.PictureDimensions()
    return float64(pictureWidth) / float64(pictureHeight)
}

// Provide the min and max plane coordinates, respectively, as defined by the user
func (app *GodelbrotApp) BigUserCoords() (BigComplex, BigComplex) {
    info := app.info
    min := BigComplex{info.BigRealMin, info.BigImagMin}
    max := BigComplex{info.BigRealMax, info.BigImagMax}
    return min, max
}

// Provide the min and max plane coordinates, respectively, as defined by the user
func (app *GodelbrotApp) NativeUserCoords() (complex128, complex128) {
    info := app.info
    return complex(info.RealMin, info.ImagMin), complex(info.RealMax, info.ImagMax)
}

func (app *GodelbrotApp) FixAspect() bool {
    return app.info.UserDescription.FixAspect
}

func (app *GodelbrotApp) SequentialNumerics() SequentialNumerics {
    return app.numerics
}

func (app *GodelbrotApp) RegionNumerics() RegionNumerics {
    return app.numerics
}

func (app *GodelbrotApp) RegionConfig() RegionParameters {
    desc := app.info.UserDescription
    glitchSamples := desc.GlitchSamples
    collapse := desc.RegionCollapse
    return RegionParameters{
        GlitchSamples: glitchSamples,
        RegionCollapseSize: collapse,
    }
}

func (app *GodelbrotApp) ConcurrentConfig() ConcurrentRegionParameters {
    desc := app.info.UserDescription
    bufferSize := desc.ThreadBufferSize
    jobs := desc.Jobs
    return ConcurrentRegionParameters{
        Jobs: jobs,
        BufferSize: bufferSize,
    }
}

func (app *GodelbrotApp) Draw() DrawingContext {
    return appDrawingContext(app)
}

// A drawing app based on the GodelbrotApp
type appDrawingContext *GodelbrotApp

func (draw appDrawingContext) Paint() Palette {
    return draw.palette
}

func (draw appDrawingContext) Picture() *image.NRGBA {
    return draw.picture
}