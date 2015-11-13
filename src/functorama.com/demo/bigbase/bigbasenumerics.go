package bigbase

import (
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

	Runit big.Float
	Iunit big.Float

	Precision uint
}

func CreateBigBaseNumerics(app RenderApplication) BigBaseNumerics {
	prec := app.Precision()

	planeMin, planeMax := app.BigUserCoords()

	left := CreateBigFloat(0.0, prec)
	right := CreateBigFloat(0.0, prec)
	top := CreateBigFloat(0.0, prec)
	bottom := CreateBigFloat(0.0, prec)

	left.Set(planeMin.Real())
	right.Set(planeMax.Real())
	bottom.Set(planeMin.Imag())
	top.Set(planeMax.Imag())

	planeWidth := CreateBigFloat(0.0, prec)
	planeWidth.Sub(&right, &left)

	planeHeight := CreateBigFloat(0.0, prec)
	planeHeight.Sub(&top, &bottom)

	planeAspect := CreateBigFloat(0.0, prec)
	planeAspect.Quo(&planeWidth, &planeHeight)

	nativePictureAspect := base.AppPictureAspectRatio(app)
	pictureAspect := CreateBigFloat(nativePictureAspect, prec)

	thindicator := planeAspect.Cmp(&pictureAspect)

	baseConfig := app.BaseConfig()

	if baseConfig.FixAspect {
		// If the plane aspect is greater than image aspect
		// Then the plane is too short, so must be made taller
		if thindicator == 1 {
			taller := CreateBigFloat(0.0, prec)
			taller.Quo(&planeWidth, &pictureAspect)
			bottom.Sub(&top, &taller)
			planeWidth.Sub(&top, &bottom)
		} else if thindicator == -1 {
			// If the plane aspect is less than the image aspect
			// Then the plane is too thin, and must be made fatter
			fatter := CreateBigFloat(0.0, prec)
			fatter.Mul(&planeHeight, &pictureAspect)
			right.Add(&left, &fatter)
			planeHeight.Sub(&right, &left)
		}
	}

	pictureWidth, pictureHeight := app.PictureDimensions()
	uq := UnitQuery{pictureWidth, pictureHeight, &planeWidth, &planeHeight}
	rUnit, iUnit := uq.PixelUnits()

	fSqrtDiverge := math.Sqrt(baseConfig.DivergeLimit)

	bbn := BigBaseNumerics{
		BaseNumerics: base.CreateBaseNumerics(app),
		RealMin:      left,
		RealMax:      right,
		ImagMin:      bottom,
		ImagMax:      top,

		SqrtDivergeLimit: CreateBigFloat(fSqrtDiverge, prec),

		Runit:     rUnit,
		Iunit:     iUnit,
		Precision: prec,
	}

	return bbn
}
func (bbn *BigBaseNumerics) CreateBigFloat(x float64) big.Float {
	return CreateBigFloat(x, bbn.Precision)
}

func (bbn *BigBaseNumerics) CreateBigComplex(r, i float64) BigComplex {
	return BigComplex{bbn.CreateBigFloat(r), bbn.CreateBigFloat(i)}
}

func (bbn *BigBaseNumerics) PlaneToPixel(c *BigComplex) (rx int, ry int) {
	// Translate x
	x := bbn.CreateBigFloat(0.0)
	x.Sub(c.Real(), &bbn.RealMin)
	// Scale x
	x.Quo(&x, &bbn.Runit)

	// Translate y
	y := bbn.CreateBigFloat(0.0)
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

func (bbn *BigBaseNumerics) CreateMandelbrotMember(c *BigComplex) BigMandelbrotMember {
	return BigMandelbrotMember{
		C: c,
		Prec: bbn.Precision,
		SqrtDivergeLimit: &bbn.SqrtDivergeLimit,
	}
}

type UnitQuery struct {
	pictureW uint
	pictureH uint
	planeW *big.Float
	planeH *big.Float
}

func (uq UnitQuery) PixelUnits() (big.Float, big.Float) {
	prec := uq.planeW.Prec()

	bigPicWidth := CreateBigFloat(float64(uq.pictureW), prec)
	bigPicHeight := CreateBigFloat(float64(uq.pictureH), prec)

	rUnit := CreateBigFloat(0.0, prec)
	rUnit.Quo(uq.planeW, &bigPicWidth)
	iUnit := CreateBigFloat(0.0, prec)
	iUnit.Quo(uq.planeH, &bigPicHeight)

	return rUnit, iUnit
}