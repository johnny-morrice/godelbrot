package libgodelbrot

import (
    "math/big"
    "image"
    "functorama.com/demo/base"
)

// Info completely describes the render process
type Info struct {
    UserDescription Request
    // Describe the render strategy in use
    DetectedRenderStrategy RenderMode
    // Describe the numerics system in use
    DetectedNumericsMode NumericsMode
    Precision int
    RealMin big.Float
    RealMax big.Float
    BigImagMin big.Float
    BigImagMax big.Float
}

// Object to initialize the godelbrot system
type configurator Info

// InitializeContext examines the description, chooses a renderer, numerical system and palette.
func configure(user *Request) (*Info, error) {
    anything, err := panic2err(func() interface{} {
        c := &configurator(Info{})
        c.UserDescription = *desc

        c.chooseNumerics()
        c.chooseRenderStrategy()

        return c
    })

    if err == nil {
        return anything.(*Info)
    } else {
        return nil, err
    }
}

// Initialize the render system
func (c *configurator) chooseRenderStrategy() {
    desc := c.UserDescription
    switch desc.RenderMode {
    case AutoDetectRenderMode:
        c.chooseFastRenderStrategy()
    case SequentialRenderMode:
        c.useSequentialRenderer()
    case RegionRenderMode:
        c.useRegionRenderer()
    case ConcurrentRegionRenderMode:
        c.useConcurrentRegionRenderer()
    default:
        log.Panic("Unknown render mode:", desc.RenderMode)
    }
}

// Initialize the numerics system
func (c *configurator) chooseNumerics() {
    desc := c.UserDescription
    c.parseUserCoords()
    switch desc.Numerics {
    case AutoDetectNumericsMode:
        c.chooseAccurateNumerics()
    case NativeNumericsMode:
        c.useNative()
    case BigFloatNumericsMode:
        c.useBig()
    default:
        log.Panic("Unknown numerics mode:", desc.Numerics)
    }
}


func (c *configurator) chooseAccurateNumerics() {
    // 53 bits precision is available to 64 bit floats
    const prec64 = 53

    desc := c.UserDescription

    prec := c.howManyBits()
    if prec > prec64 {
        c.useBig()
        c.setPrec(prec)
    } else {
        c.useNative()
        c.setPrec(prec64)
    }
}

func (c *configurator) setPrec(bits int) {
    c.Precision = bits
    bounds := []*big.Float{
        &c.RealMin,
        &c.RealMax,
        &c.ImagMin,
        &c.ImagMax,
    }

    for num := range bounds {
        // I say c.Precision rather than bits because I think these should be equal
        // and if there is a bug, this will certainly break quicker.
        num.SetPrec(c.Precision)
    }
}

func (c *configurator) useNative() {
    c.DetectedNumericsMode = NativeNumericsMode
}

func (c *configurator) useBig() {
    c.DetectedNumericsMode = BigFloatNumericsMode
}

func (c *configurator) parseUserCoords() {
    bigActions := []func(big.Float){
        func(realMin big.Float) { c.BigRealMin = realMin },
        func(realMax big.Float) { c.BigRealMax = realMax },
        func(imagMin big.Float) { c.BigImagMin = imagMin },
        func(imagMax big.Float) { c.BigImagMax = imagMax },
    }

    desc := c.UserDescription
    userInput := []string{
        desc.RealMin,
        desc.RealMax,
        desc.ImagMin,
        desc.ImagMax,
    }

    inputNames := []string{"realMin", "realMax", "imagMin", "imagMax"}

    for i, num := range userInput {
        bigFloat, bigErr := parseBig(num)

        badName := inputNames[i]

        if bigErr != nil {
            parsePanic(bigErr, badName)
        }

        // Handle parse results
        nativeActions[i](native)
        bigActions[i](bigFloat)
    }
}

// Choose an optimal strategy for rendering the image
func (c *configurator) chooseFastRenderStrategy() {
    desc := c.UserDescription

    area := desc.ImageWidth * desc.ImageHeight
    numerics := c.DetectedNumericsMode

    if numerics == AutoDetectNumericsMode {
        log.Panic("Must choose render strategy after numerics system")
    }

    if area < DefaultTinyImageArea && numerics == NativeNumericsMode {
        // Use `SequenceRenderStrategy' when
        // We have native arithmetic and the image is tiny
        c.useSequentialRenderer()
    } else if desc.Jobs <= DefaultLowThreading {
        // Use `RegionRenderStrategy' when
        // the number of jobs is small
        c.useRegionRenderer()
    } else {
        // Use `ConcurrentRegionRenderStrategy' otherwise
        c.useConcurrentRegionRenderer()
    }
}

func (c *configurator) useSequentialRenderer() {
    c.DetectedRenderStrategy = SequenceRenderMode
}

func (c *configurator) useRegionRenderer() {
    c.DetectedRenderStrategy = RegionRenderMode
}

func (c *configurator) useConcurrentRegionRenderer() {
    c.DetectedRenderStrategy = ConcurrentRegionRenderMode
}

func (c *configurator) howManyBits() int {
    desc := &c.UserDescription
    bottom, top, imageLength := c.mostSqueezed()

    divisions := desc.IterLimit * imageLength
    // I wonder how long this will take!
    for bits := 1; !space(divisions, bits, &bottom, &top); bits++ {
        // A location for some progress counter hook?
    }

    return bits
}

func (c *configurator) mostSqueezed() (big.Float, big.Float int) {
    desc := &c.UserDescription

    picWidth := big.Float{}
    picWidth.SetFloat64(float64(desc.ImageWidth))

    picHeight := big.Float{}
    picHeight.SetFloat64(float64(desc.ImageHeight))

    rLength, rShrink := axisShrink(&c.RealMax, &c.RealMin, &picWidth)
    iLength, iShrink := axisShrink(&c.ImagMax, &c.ImagMin, &picHeight)

    if rShrink < iShrink {
        return desc.RealMin, desc.RealMax, desc.ImageWidth
    } else {
        return desc.ImagMin, desc.ImagMax, desc.ImageHeight
    }
}

func axisShrink(max *big.Float, min *big.Float, outLenght *big.Float) (big.Float, big.Float) {
    length := big.Float{}
    length.Copy(max)
    length.Sub(&length, min)
    shrink := big.Float{}
    shrink.Copy(&rLength)
    shrink.Quo(&shrink, &picWidth)
    return shrink, length
}

func space(divisions int, bits int, bottom *big.Float, top *big.Float) bool {
    // Copy bottom pointer
    var current *big.Float = &big.Float{}
    *current = *bottom

    current.SetPrec(bits)
    for i := 0; i < int(divisions); i++ {
        current = nextafter(&current, top)
    }
    // True is current less than top
    return current.Cmp(top) == -1
}

func nextafter(from *big.Float, to *big.Float) *big.Float {
}