package nativebase

import (
	"math"
	"functorama.com/demo/base"
)

// Basis for all native numerics
type NativeBaseNumerics struct {
	base.BaseNumerics

	RealMin float64
	RealMax float64
	ImagMin float64
	ImagMax float64

	Runit float64
	Iunit float64

	DivergeLimit float64
}

func CreateNativeBaseNumerics(app NativeRenderApplication) NativeBaseNumerics {
	planeMin, planeMax := app.NativeUserCoords()
	planeWidth := real(planeMax) - real(planeMin)
	planeHeight := imag(planeMax) - imag(planeMin)
	planeAspect := planeWidth / planeHeight
	pictureWidth, pictureHeight := app.PictureDimensions()
	pictureAspect := base.PictureAspectRatio(pictureWidth, pictureHeight)

	if app.FixAspect() {
		// If the plane aspect is greater than image aspect
		// Then the plane is too short, so must be made taller
		if planeAspect > pictureAspect {
			taller := planeWidth / pictureAspect
			bottom := imag(planeMax) - taller
			planeMin = complex(real(planeMin), bottom)
		} else if planeAspect < pictureAspect {
			// If the plane aspect is less than the image aspect
			// Then the plane is too thin, and must be made fatter
			fatter := planeHeight * pictureAspect
			right := real(planeMax) + fatter
			planeMax = complex(right, imag(planeMax))
		}
	}

	limit := app.DivergeLimit()

	return NativeBaseNumerics{
		BaseNumerics: base.CreateBaseNumerics(app),
		RealMin:      real(planeMin),
		RealMax:      real(planeMax),
		ImagMin:      imag(planeMin),
		ImagMax:      imag(planeMax),

		DivergeLimit: limit,

		Runit: planeWidth / float64(pictureWidth),
		Iunit: planeHeight / float64(pictureHeight),
	}
}

func (native NativeBaseNumerics) PlaneTopLeft() complex128 {
	return complex(native.RealMin, native.ImagMax)
}

// Size on the plane of 1px
func (native NativeBaseNumerics) PixelSize() (float64, float64) {
	return native.Runit, native.Iunit
}

func (native NativeBaseNumerics) PlaneToPixel(c complex128) (rx int, ry int) {
	topLeft := native.PlaneTopLeft()
	rUnit, iUnit := native.PixelSize()
	// Translate x
	tx := real(c) - real(topLeft)
	// Scale x
	sx := tx / rUnit

	// Translate y
	ty := imag(c) - imag(topLeft)
	// Scale y
	sy := ty / iUnit

	rx = int(math.Floor(sx))
	// Remember that we draw downwards
	ry = int(math.Ceil(-sy))

	return
}
