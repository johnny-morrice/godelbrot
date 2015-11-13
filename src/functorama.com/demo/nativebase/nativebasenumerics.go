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

	SqrtDivergeLimit float64
}

func CreateNativeBaseNumerics(app RenderApplication) NativeBaseNumerics {
	planeMin, planeMax := app.NativeUserCoords()
	planeWidth := real(planeMax) - real(planeMin)
	planeHeight := imag(planeMax) - imag(planeMin)
	planeAspect := planeWidth / planeHeight
	pictureWidth, pictureHeight := app.PictureDimensions()
	pictureAspect := base.PictureAspectRatio(pictureWidth, pictureHeight)
	config := app.BaseConfig()

	if config.FixAspect {
		// If the plane aspect is greater than image aspect
		// Then the plane is too short, so must be made taller
		if planeAspect > pictureAspect {
			taller := planeWidth / pictureAspect
			bottom := imag(planeMax) - taller
			planeMin = complex(real(planeMin), bottom)
			planeHeight = imag(planeMax) - bottom
		} else if planeAspect < pictureAspect {
			// If the plane aspect is less than the image aspect
			// Then the plane is too thin, and must be made fatter
			fatter := planeHeight * pictureAspect
			right := real(planeMax) + fatter
			planeMax = complex(right, imag(planeMax))
			planeWidth = right - real(planeMin)
		}
	}

	uq := UnitQuery{pictureWidth, pictureHeight, planeWidth, planeHeight}
	rUnit, iUnit := uq.PixelUnits()

	return NativeBaseNumerics{
		BaseNumerics: base.CreateBaseNumerics(app),
		RealMin:      real(planeMin),
		RealMax:      real(planeMax),
		ImagMin:      imag(planeMin),
		ImagMax:      imag(planeMax),

		SqrtDivergeLimit: math.Sqrt(config.DivergeLimit),

		Runit: rUnit,
		Iunit: iUnit,
	}
}

func (native *NativeBaseNumerics) CreateMandelbrot(c complex128) NativeMandelbrotMember {
	return NativeMandelbrotMember{
		C: c,
		SqrtDivergeLimit: native.SqrtDivergeLimit,
	}
}

func (native *NativeBaseNumerics) PlaneTopLeft() complex128 {
	return complex(native.RealMin, native.ImagMax)
}

// Size on the plane of 1px
func (native *NativeBaseNumerics) PixelSize() (float64, float64) {
	return native.Runit, native.Iunit
}

func (native *NativeBaseNumerics) PlaneToPixel(c complex128) (rx int, ry int) {
	const debug = true
	topLeft := native.PlaneTopLeft()
	rUnit, iUnit := native.PixelSize()
	// Translate x
	tx := real(c) - real(topLeft)

	// Translate y
	ty := imag(topLeft) - imag(c)

	// Scale x
	sx := tx / rUnit
	// Scale y
	sy := ty / iUnit

	rx = int(math.Floor(sx))
	// Remember that we draw downwards
	ry = int(math.Ceil(sy))

	return
}

type UnitQuery struct {
	pictureW uint
	pictureH uint
	planeW float64
	planeH float64
}

func (uq UnitQuery) PixelUnits() (float64, float64) {
	rUnit := uq.planeW / float64(uq.pictureW)
	iUnit := uq.planeH / float64(uq.pictureH)
	return rUnit, iUnit
}