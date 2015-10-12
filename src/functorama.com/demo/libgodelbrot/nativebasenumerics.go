package libgodelbrot

type NativeBaseNumerics struct {
    realMin float64
    realMax float64
    imagMin float64
    imagMax float64

    picXMin int
    picXMax int
    picYMin int
    picYMax int

    collapse int

    iterLimit int
    divergeLimit int

    rUnit float64
    iUnit float64
}

func NewNativeBaseNumerics(context *ContextFacade) *NativeBaseNumerics {
    planeMin, planeMax := context.NativeUserCoords()
    planeWidth := real(planeMax) - real(planeMin)
    planeHeight := imag(planeMax) - imag(planeMin)
    pictureWidth, pictureHeight := context.PictureDimensions()
    planeAspect := planeWidth / planeHeight
    pictureAspect := float64(pictureWidth) / float64(pictureHeight)

    if context.FixAspect() {
        // If the plane aspect is greater than image aspect
        // Then the plane is too short, so must be made taller
        if planeAspect > pictureAspect {
            taller := planeWidth / pictureAspect
            bottom := imag(planeMax) - taller
            planeMin = complex128(real(planeMin, bottom))
        } else if planeAspect < pictureAspect {
            // If the plane aspect is less than the image aspect
            // Then the plane is too thin, and must be made fatter
            fatter := planeHeight * pictureAspect
            right := real(planeMax) + fatter
            planeMax = complex128(right, imag(planeMax))
        }
    }

    iLimit, dLimit := context.Limits()

    return &NativeBaseNumerics{
        realMin: real(planeMin),
        realMax: real(planeMax),
        imagMin: imag(planeMin),
        imagMax: imag(planeMax),

        picXMin: 0,
        picXMax: pictureWidth,
        picYMin: 0,
        picYMax: pictureHeight,

        collapse: context.RegionCollapseSize(),

        iterLimit: iLimit,
        divergeLimit: dLimit,

        rUnit: planeWidth / float64(pictureWidth),
        iUnit: planeHeight / float64(pictureHeight),
    }
}

func (native NativeBaseNumerics) PictureMin() (int, int) {
    return native.picXMin, native.picYMin
}

func (native NativeBaseNumerics) PictureMax() (int, int) {
    return native.picXMax, native.picYMax
}

func (native NativeBaseNumerics) PlaneTopLeft() complex128 {
    return complex(native.realMin, native.imagMax)
}

// Size on the plane of 1px
func (native NativeBaseNumerics) PixelSize() (float64, float64) {
    return rUnit, iUnit
}

func (native NativeBaseNumerics) IterateLimit() int {
    return native.iterLimit
}

func (native NativeBaseNumerics) DivergeLimit() float64 {
    return native.divergeLimit
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

    rx = math.Floor(sx)
    // Remember that we draw downwards
    ry = math.Ceil(-sy)

    return
}
