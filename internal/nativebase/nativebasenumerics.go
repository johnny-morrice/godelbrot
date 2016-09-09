package nativebase

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"image"
	"math"
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
	IterateLimit     uint8
}

func Make(app RenderApplication) NativeBaseNumerics {
	planeMin, planeMax := app.NativeUserCoords()
	planeWidth := real(planeMax) - real(planeMin)
	planeHeight := imag(planeMax) - imag(planeMin)
	pictureWidth, pictureHeight := app.PictureDimensions()
	config := app.BaseConfig()

	uq := UnitQuery{pictureWidth, pictureHeight, planeWidth, planeHeight}
	rUnit, iUnit := uq.PixelUnits()

	return NativeBaseNumerics{
		BaseNumerics: base.Make(app),
		RealMin:      real(planeMin),
		RealMax:      real(planeMax),
		ImagMin:      imag(planeMin),
		ImagMax:      imag(planeMax),

		SqrtDivergeLimit: math.Sqrt(config.DivergeLimit),
		IterateLimit:     config.IterateLimit,

		Runit: rUnit,
		Iunit: iUnit,
	}
}

func (nbn *NativeBaseNumerics) CreateMandelbrot(c complex128) NativeEscapeValue {
	return NativeEscapeValue{
		C:                c,
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

	rx = round(sx)
	// Remember that we draw downwards
	ry = round(sy)

	return
}

func (nbn *NativeBaseNumerics) PixelToPlane(i, j int) complex128 {
	tr := nbn.Xtor(i)
	ti := nbn.Ytoi(j)

	return complex(tr, ti)
}

func (nbn *NativeBaseNumerics) Xtor(i int) float64 {
	sr := float64(i) * nbn.Runit

	return nbn.RealMin + sr
}

func (nbn *NativeBaseNumerics) Ytoi(j int) float64 {
	si := float64(j) * nbn.Iunit

	return nbn.ImagMax - si
}

func (nbn *NativeBaseNumerics) Escape(c complex128) NativeEscapeValue {
	point := nbn.CreateMandelbrot(c)
	point.Mandelbrot(nbn.IterateLimit)
	return point
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
	PlaneW   float64
	PlaneH   float64
}

func (uq UnitQuery) PixelUnits() (float64, float64) {
	rUnit := uq.PlaneW / float64(uq.PictureW)
	iUnit := uq.PlaneH / float64(uq.PictureH)
	return rUnit, iUnit
}

func round(r float64) int {
	frac := math.Abs(r - math.Floor(r))
	if frac >= 0.5 {
		return int(math.Ceil(r))
	} else {
		return int(math.Floor(r))
	}
}
