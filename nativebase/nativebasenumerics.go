package nativebase

import (
	"image"
	"math"
	"github.com/johnny-morrice/godelbrot/base"
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
	IterateLimit uint8
}

func Make(app RenderApplication) NativeBaseNumerics {
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
			right := real(planeMin) + fatter
			planeMax = complex(right, imag(planeMax))
			planeWidth = right - real(planeMin)
		}
	}

	uq := UnitQuery{pictureWidth, pictureHeight, planeWidth, planeHeight}
	rUnit, iUnit := uq.PixelUnits()

	return NativeBaseNumerics{
		BaseNumerics: base.Make(app),
		RealMin:      real(planeMin),
		RealMax:      real(planeMax),
		ImagMin:      imag(planeMin),
		ImagMax:      imag(planeMax),

		SqrtDivergeLimit: math.Sqrt(config.DivergeLimit),
		IterateLimit: config.IterateLimit,

		Runit: rUnit,
		Iunit: iUnit,
	}
}

func (nbn *NativeBaseNumerics) CreateMandelbrot(c complex128) NativeEscapeValue {
	return NativeEscapeValue{
		C: c,
		SqrtDivergeLimit: nbn.SqrtDivergeLimit,
	}
}

// Size on the plane of 1px
func (nbn *NativeBaseNumerics) PixelSize() (float64, float64) {
	return nbn.Runit, nbn.Iunit
}

func (nbn *NativeBaseNumerics) PlaneToPixel(c complex128) (rx int, ry int) {
	rUnit, iUnit := nbn.PixelSize()
	// Translate x
	tx := real(c) - nbn.RealMin

	// Translate y
	ty := nbn.ImagMax - imag(c)

	// Scale x
	sx := tx / rUnit
	// Scale y
	sy := ty / iUnit

	rx = int(math.Floor(sx))
	// Remember that we draw downwards
	ry = int(math.Ceil(sy))

	return
}

func (nbn *NativeBaseNumerics) PixelToPlane(i, j int) complex128 {
	rUnit, iUnit := nbn.PixelSize()

	// Scale
	sr := float64(i) * rUnit
	si := float64(j) * iUnit

	// Translate
	tr := nbn.RealMin + sr
	ti := nbn.ImagMax - si

	return complex(tr, ti)
}

func (nbn *NativeBaseNumerics) SubImage(rect image.Rectangle) {
	min := nbn.PixelToPlane(rect.Min.X, rect.Min.Y)
	max := nbn.PixelToPlane(rect.Max.X, rect.Max.Y)

	nbn.PictureSubImage(rect)

	nbn.RealMin = real(min)
	nbn.ImagMin = imag(min)
	nbn.RealMax = real(max)
	nbn.ImagMax = imag(max)
}

type UnitQuery struct {
	PictureW uint
	PictureH uint
	PlaneW float64
	PlaneH float64
}

func (uq UnitQuery) PixelUnits() (float64, float64) {
	rUnit := uq.PlaneW / float64(uq.PictureW)
	iUnit := uq.PlaneH / float64(uq.PictureH)
	return rUnit, iUnit
}