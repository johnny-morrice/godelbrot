package bigbase

import (
	"image"
	"math"
	"math/big"
	"functorama.com/demo/base"
)

// Basis for all big.Float numerics
type BigBaseNumerics struct {
	base.BaseNumerics

	RealMin big.Float
	RealMax big.Float
	ImagMin big.Float
	ImagMax big.Float

	SqrtDivergeLimit big.Float
	IterateLimit uint8

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

	planeAspect := MakeBigFloat(0.0, prec)
	planeAspect.Quo(&planeWidth, &planeHeight)

	nativePictureAspect := base.AppPictureAspectRatio(app)
	pictureAspect := MakeBigFloat(nativePictureAspect, prec)

	thindicator := planeAspect.Cmp(&pictureAspect)

	baseConfig := app.BaseConfig()

	if baseConfig.FixAspect {
		// If the plane aspect is greater than image aspect
		// Then the plane is too short, so must be made taller
		if thindicator == 1 {
			taller := MakeBigFloat(0.0, prec)
			taller.Quo(&planeWidth, &pictureAspect)
			bottom.Sub(&top, &taller)
			planeHeight.Sub(&top, &bottom)
		} else if thindicator == -1 {
			// If the plane aspect is less than the image aspect
			// Then the plane is too thin, and must be made fatter
			fatter := MakeBigFloat(0.0, prec)
			fatter.Mul(&planeHeight, &pictureAspect)
			right.Add(&left, &fatter)
			planeWidth.Sub(&right, &left)
		}
	}

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
		IterateLimit: baseConfig.IterateLimit,

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

func (bbn *BigBaseNumerics) MakeMember(c *BigComplex) BigMandelbrotMember {
	return BigMandelbrotMember{
		C: c,
		Prec: bbn.Precision,
		SqrtDivergeLimit: &bbn.SqrtDivergeLimit,
	}
}

// TODO
func (bbn *BigBaseNumerics) SubImage(rect image.Rectangle) {

}

type UnitQuery struct {
	pictureW uint
	pictureH uint
	planeW *big.Float
	planeH *big.Float
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