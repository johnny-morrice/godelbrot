package libgodelbrot

import (
    "math/big"
)

// Basis for all big.Float numerics
type BigBaseNumerics struct {
    BaseNumerics
    BaseRegionNumerics

    realMin big.Float
    realMax big.Float
    imagMin big.Float
    imagMax big.Float

    divergeLimit big.Float

    rUnit big.Float
    iUnit big.Float

    precision uint
}

func CreateBigBaseNumerics(context *ContextFacade) BigBaseNumerics {
    prec := pixelPerfectPrecision(context)

    planeMin, planeMax := context.BigUserCoords()

    left := NewBigFloat(0.0, prec)
    right := NewBigFloat(0.0, prec)
    bottom := NewBigFloat(0.0, prec)
    top := NewBigFloat(0.0, prec)

    left.Set(planeMin.Real())
    right.Set(planeMax.Real())
    bottom.Set(planeMin.Imag()    
    top.Set(planeMax.Imag()))

    planeWidth := NewBigFloat(0.0, prec)
    planeWidth.Sub(right, left)

    planeHeight := NewBigFloat(0.0, prec)
    planeHeight.Sub(top, bottom)

    planeAspect := NewBigFloat(0.0, prec)
    planeAspect.Quo(planeWidth, planeHeight)

    nativePictureAspect := context.PictureAspect()
    pictureAspect := NewBigFloat(nativePictureAspect, prec)

    thindicator := planeAspect.Cmp(pictureAspect)

    if context.FixAspect() {
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

    iLimit, dLimit := context.Limits()

    pictureWidthI, pictureHeightI := context.PictureDimensions()
    pictureWidth := NewBigFloat(float64(pictureWidthI), prec)
    pictureHeight := NewBigFloat(float64(pictureHeightI), prec)

    rSize := NewBigFloat(0.0, prec)
    rSize.Quo(planeWidth, pictureWidth)
    iSize := NewBigFloat(0.0, prec)
    iSize.Quo(planeHeight, pictureHeight)

    return BigBaseNumerics{
        BaseNumerics: CreateBaseNumerics(context),
        realMin: real(planeMin),
        realMax: real(planeMax),
        imagMin: imag(planeMin),
        imagMax: imag(planeMax),

        divergeLimit: NewBigFloat(dLimit, prec),

        rUnit: rSize,
        iUnit: iSize,
        precision: prec,
    }
}

func (bigFloat BigBaseNumerics) NewBigFloat(f float64) big.Float {
    return NewBigFloat(f, big.precision)
}

func (bigFloat BigBaseNumerics) PictureMin() (int, int) {
    return bigFloat.picXMin, bigFloat.picYMin
}

func (bigFloat BigBaseNumerics) PictureMax() (int, int) {
    return bigFloat.picXMax, bigFloat.picYMax
}

func (bigFloat BigBaseNumerics) PlaneTopLeft() BigComplex {
    return BigComplex{bigFloat.realMin, bigFloat.imagMax}
}

// Size on the plane of 1px
func (bigFloat BigBaseNumerics) PixelSize() (big.Float, big.Float) {
    return rUnit, iUnit
}

func (bigFloat BigBaseNumerics) MandelbrotLimits (int, big.Float) {
    return bigFloat.iterLimit, bigFloat.divergeLimit
}

func (bigFloat BigBaseNumerics) PlaneToPixel(c BigComplex) (rx int, ry int) {
    topLeft := bigFloat.PlaneTopLeft()
    rUnit, iUnit := bigFloat.PixelSize()

    // Translate x
    x := bigFloat.NewBigFloat(0.0)
    x.Sub(c.Real(), topLeft.Real())
    // Scale x
    x.Quo(x, rUnit)

    // Translate y
    y := bigFloat.NewBigFloat(0.0)
    y.Sub(c.Imag(), topLeft.Imag())
    // Scale y
    y.Quo(y, iUnit)

    fx, _ = x.Float64()
    fy, _ = y.Float64()

    rx = math.Floor(fx)
    // Remember that we draw downwards
    ry = math.Ceil(-fy)

    return
}
