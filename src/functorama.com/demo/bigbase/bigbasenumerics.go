package bigbase

import (
	"math"
	"math/big"
)

// Basis for all big.Float numerics
type BigBaseNumerics struct {
	BaseNumerics
	BaseRegionNumerics

	RealMin big.Float
	RealMax big.Float
	ImagMin big.Float
	ImagMax big.Float

	DivergeLimit big.Float

	Runit big.Float
	Iunit big.Float

	Precision uint
}

func CreateBigBaseNumerics(app RenderApplication) BigBaseNumerics {
	prec := DefaultHighPrec

	planeMin, planeMax := app.BigUserCoords()

	left := NewBigFloat(0.0, prec)
	right := NewBigFloat(0.0, prec)
	bottom := NewBigFloat(0.0, prec)
	top := NewBigFloat(0.0, prec)

	left.Set(planeMin.Real())
	right.Set(planeMax.Real())
	bottom.Set(planeMin.Imag())
	top.Set(planeMax.Imag())

	planeWidth := NewBigFloat(0.0, prec)
	planeWidth.Sub(right, left)

	planeHeight := NewBigFloat(0.0, prec)
	planeHeight.Sub(top, bottom)

	planeAspect := NewBigFloat(0.0, prec)
	planeAspect.Quo(planeWidth, planeHeight)

	nativePictureAspect := AppPictureAspectRatio(app)
	pictureAspect := NewBigFloat(nativePictureAspect, prec)

	thindicator := planeAspect.Cmp(pictureAspect)

	if app.FixAspect() {
		// If the plane aspect is greater than image aspect
		// Then the plane is too short, so must be made taller
		if thindicator == 1 {
			taller := NewBigFloat(0.0, prec)
			taller.Quo(planeWidth, pictureAspect)
			bottom.Sub(top - taller)
			planeMin = BigComplex{left, bottom}
		} else if thindicator == -1 {
			// If the plane aspect is less than the image aspect
			// Then the plane is too thin, and must be made fatter
			fatter := NewBigFloat(0.0, prec)
			fatter.Mul(planeHeight, pictureAspect)
			right.Add(left, fatter)
			planeMax = BigComplex{right, top}
		}
	}

	iLimit, dLimit := app.Limits()

	pictureWidthI, pictureHeightI := app.PictureDimensions()
	pictureWidth := NewBigFloat(float64(pictureWidthI), prec)
	pictureHeight := NewBigFloat(float64(pictureHeightI), prec)

	rSize := NewBigFloat(0.0, prec)
	rSize.Quo(planeWidth, pictureWidth)
	iSize := NewBigFloat(0.0, prec)
	iSize.Quo(planeHeight, pictureHeight)

	base := BigBaseNumerics{
		BaseNumerics: CreateBaseNumerics(app),
		RealMin:      real(planeMin),
		RealMax:      real(planeMax),
		ImagMin:      imag(planeMin),
		ImagMax:      imag(planeMax),

		DivergeLimit: NewBigFloat(dLimit, prec),

		Runit:     rSize,
		Iunit:     iSize,
		Precision: prec,
	}

	// Reduce the Precision of the base for swifter rendering
	base.FastPixelPerfectPrecision()

	return base
}

func (base *BigBaseNumerics) PictureMin() (int, int) {
	return base.picXMin, base.picYMin
}

func (base *BigBaseNumerics) PictureMax() (int, int) {
	return base.picXMax, base.picYMax
}

func (base *BigBaseNumerics) PlaneTopLeft() BigComplex {
	return BigComplex{base.RealMin, base.ImagMax}
}

// Size on the plane of 1px
func (base *BigBaseNumerics) PixelSize() (big.Float, big.Float) {
	return Runit, Iunit
}

func (base *BigBaseNumerics) MandelbrotLimits(int, big.Float) {
	return base.iterLimit, base.DivergeLimit
}

func (base *BigBaseNumerics) PlaneToPixel(c BigComplex) (rx int, ry int) {
	topLeft := base.PlaneTopLeft()
	Runit, Iunit := base.PixelSize()

	// Translate x
	x := base.NewBigFloat(0.0)
	&x.Sub(c.Real(), topLeft.Real())
	// Scale x
	&x.Quo(x, Runit)

	// Translate y
	y := base.NewBigFloat(0.0)
	y.Sub(c.Imag(), topLeft.Imag())
	// Scale y
	y.Quo(y, Iunit)

	fx, _ = x.Float64()
	fy, _ = y.Float64()

	rx = math.Floor(fx)
	// Remember that we draw downwards
	ry = math.Ceil(-fy)

	return
}

// FastPixelPerfectPrecision reduces Precision of the numeric system, while maintaining adequate
// accuracy.   Returns the new precison.
func (base *BigBaseNumerics) FastPixelPerfectPrecision() uint {
	// To keep things speedy, we will only explore 2 paths through the image
	xMin, yMin := base.PictureMin()
	xMax, yMax := base.PictureMax()

	highPrec = 0
	Runit, Iunit := base.PixelSize()

	topLeft := base.PlaneTopLeft()
	row := topLeft.Real().Copy()
	column := topLeft.Imag().Copy()

	// Find lowest required prec in the real axis
	for i := xMin; i < xMax; i++ {
		rowPrec := row.MinPrec()
		if rowPrec > highPrec {
			highPrec = rowPrec
		}
		row.Add(row, Runit)
	}

	// Find lowest required prec in the y axis
	for i := yMin; i < yMax; i++ {
		colPrec := col.MinPrec()
		if colPrec > highPrec {
			highPrec = colPrec
		}
		row.Sub(column, Iunit)
	}

	base.SetPrec(highPrec)

	return highPrec
}

// Set the Precision of the base
func (base *BigBaseNumerics) SetPrec(prec uint) {
	base.Precision = prec
	baseFloats := []big.Float{
		base.RealMin,
		base.RealMax,
		base.ImagMin,
		base.ImagMax,
		base.DivergeLimit,
		base.Runit,
		base.Iunit,
	}

	for _, f := range baseFloats {
		f.SetPrec(prec)
	}
}
