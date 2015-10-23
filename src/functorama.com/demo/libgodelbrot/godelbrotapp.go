package libgodelbrot

// An interface to the Godelbrot system at large
type GodelbrotApp *ContextInit

func (app *GodelbrotApp) IterateLimit() uint8 {
	return app.info.UserDescription.IterateLimit
}

func (app *GodelbrotApp) DivergeLimit() float64 {
	return app.info.UserDescription.DivergeLimit
}

// Provide the image dimensions
func (app *GodelbrotApp) PictureDimensions() (uint, uint) {
	desc := app.info.UserDescription
	return desc.ImageWidth, desc.ImageHeight
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

func (app *GodelbrotApp) NumericsFactory() AbstractNumericsFactory {
	return app.numerics
}

func (app *GodelbrotApp) RegionConfig() RegionParameters {
	desc := app.info.UserDescription
	glitchSamples := desc.GlitchSamples
	collapse := desc.RegionCollapse
	return RegionParameters{
		GlitchSamples:      glitchSamples,
		RegionCollapseSize: collapse,
	}
}

func (app *GodelbrotApp) ConcurrentConfig() ConcurrentRegionParameters {
	desc := app.info.UserDescription
	bufferSize := desc.ThreadBufferSize
	jobs := desc.Jobs
	return ConcurrentRegionParameters{
		Jobs:       jobs,
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
