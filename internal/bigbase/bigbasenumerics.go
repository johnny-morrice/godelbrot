package bigbase

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"image"
	"math"
	"math/big"
)

// Basis for all big.Float numerics
type BigBaseNumerics struct {
	base.BaseNumerics

	RealMin big.Float
	RealMax big.Float
	ImagMin big.Float
	ImagMax big.Float

	SqrtDivergeLimit big.Float
	IterateLimit     uint8

	Runit big.Float
	Iunit big.Float

	Precision uint
}

func Make(app RenderApplication) BigBaseNumerics {
	prec := app.Precision()

	planeMin, planeMax := app.BigUserCoords()

	left := MakeBigFloat(0.0, prec)
	right := MakeBigFloat(0.0, prec)
	top := MakeBigFloat(0.0, prec)
	bottom := MakeBigFloat(0.0, prec)

	left.Set(planeMin.Real())
	right.Set(planeMax.Real())
	bottom.Set(planeMin.Imag())
	top.Set(planeMax.Imag())

	planeWidth := MakeBigFloat(0.0, prec)
	planeWidth.Sub(&right, &left)

	planeHeight := MakeBigFloat(0.0, prec)
	planeHeight.Sub(&top, &bottom)

	baseConfig := app.BaseConfig()

	pictureWidth, pictureHeight := app.PictureDimensions()
	uq := UnitQuery{pictureWidth, pictureHeight, &planeWidth, &planeHeight}
	rUnit, iUnit := uq.PixelUnits()

	fSqrtDiverge := math.Sqrt(baseConfig.DivergeLimit)

	bbn := BigBaseNumerics{
		BaseNumerics: base.Make(app),
		RealMin:      left,
		RealMax:      right,
		ImagMin:      bottom,
		ImagMax:      top,

		SqrtDivergeLimit: MakeBigFloat(fSqrtDiverge, prec),
		IterateLimit:     baseConfig.IterateLimit,

		Runit:     rUnit,
		Iunit:     iUnit,
		Precision: prec,
	}

	return bbn
}
func (bbn *BigBaseNumerics) MakeBigFloat(x float64) big.Float {
	return MakeBigFloat(x, bbn.Precision)
}

func (bbn *BigBaseNumerics) MakeBigComplex(r, i float64) BigComplex {
	return BigComplex{bbn.MakeBigFloat(r), bbn.MakeBigFloat(i)}
}

func (bbn *BigBaseNumerics) PlaneToPixel(c *BigComplex) (rx int, ry int) {
	// Translate x
	x := bbn.MakeBigFloat(0.0)
	x.Sub(c.Real(), &bbn.RealMin)
	// Scale x
	x.Quo(&x, &bbn.Runit)

	// Translate y
	y := bbn.MakeBigFloat(0.0)
	y.Sub(&bbn.ImagMax, c.Imag())
	// Scale y
	y.Quo(&y, &bbn.Iunit)

	fx, _ := x.Float64()
	fy, _ := y.Float64()

	rx = int(math.Floor(fx))
	// Remember that we draw downwards
	ry = int(math.Ceil(fy))

	return
}

func (bbn *BigBaseNumerics) MakeMember(c *BigComplex) BigEscapeValue {
	return BigEscapeValue{
		C:                c,
		Prec:             bbn.Precision,
		SqrtDivergeLimit: &bbn.SqrtDivergeLimit,
	}
}

func (bbn *BigBaseNumerics) SubImage(rect image.Rectangle) {
	min := bbn.PixelToPlane(rect.Min.X, rect.Min.Y)
	max := bbn.PixelToPlane(rect.Max.X, rect.Max.Y)

	bbn.PictureSubImage(rect)

	bbn.RealMin = min.R
	bbn.ImagMin = min.I
	bbn.RealMax = max.R
	bbn.ImagMax = max.I
}

func (bbn *BigBaseNumerics) PixelToPlane(i, j int) BigComplex {
	rUnit, iUnit := bbn.PixelSize()

	x := bbn.MakeBigFloat(float64(i))
	y := bbn.MakeBigFloat(float64(j))

	// Scale
	re := bbn.MakeBigFloat(0.0)
	re.Mul(&x, &rUnit)

	im := bbn.MakeBigFloat(0.0)
	im.Mul(&y, &iUnit)

	// Translate
	re.Add(&re, &bbn.RealMin)

	// Dodge aliasing error
	extra := bbn.MakeBigFloat(0.0)
	extra.Sub(&bbn.ImagMax, &im)

	return BigComplex{re, extra}
}

// Size on the plane of 1px
func (bbn *BigBaseNumerics) PixelSize() (big.Float, big.Float) {
	return bbn.Runit, bbn.Iunit
}

func (bbn *BigBaseNumerics) Escape(c *BigComplex) BigEscapeValue {
	point := bbn.MakeMember(c)
	point.Mandelbrot(bbn.IterateLimit)
	return point
}

type UnitQuery struct {
	pictureW uint
	pictureH uint
	planeW   *big.Float
	planeH   *big.Float
}

func (uq UnitQuery) PixelUnits() (big.Float, big.Float) {
	prec := uq.planeW.Prec()

	bigPicWidth := MakeBigFloat(float64(uq.pictureW), prec)
	bigPicHeight := MakeBigFloat(float64(uq.pictureH), prec)

	rUnit := MakeBigFloat(0.0, prec)
	rUnit.Quo(uq.planeW, &bigPicWidth)
	iUnit := MakeBigFloat(0.0, prec)
	iUnit.Quo(uq.planeH, &bigPicHeight)

	return rUnit, iUnit
}
