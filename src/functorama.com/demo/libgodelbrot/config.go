package libgodelbrot

import (
    "math/big"
    "log"
)

// Object to initialize the godelbrot system
type configurator Info

// InitializeContext examines the description, chooses a renderer, numerical system and palette.
func configure(req *Request) *Info {
    c := &configurator{}
    c.UserRequest = *req

    c.chooseNumerics()
    c.chooseRenderStrategy()

    return (*Info)(c)
}

// Initialize the render system
func (c *configurator) chooseRenderStrategy() {
    req := c.UserRequest
    switch req.Renderer {
    case AutoDetectRenderMode:
        c.chooseFastRenderStrategy()
    case SequenceRenderMode:
        c.useSequenceRenderer()
    case RegionRenderMode:
        c.useRegionRenderer()
    case SharedRegionRenderMode:
        c.useSharedRegionRenderer()
    default:
        log.Panic("Unknown render mode:", req.Renderer)
    }
}

// Initialize the numerics system
func (c *configurator) chooseNumerics() {
    desc := c.UserRequest
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
    const prec64 uint = 53

    prec := c.howManyBits()
    if prec > prec64 {
        c.useBig()
        c.setPrec(prec)
    } else {
        c.useNative()
        c.setPrec(prec64)
    }
}

func (c *configurator) setPrec(bits uint) {
    c.Precision = bits
    bounds := []*big.Float{
        &c.RealMin,
        &c.RealMax,
        &c.ImagMin,
        &c.ImagMax,
    }

    for _, num := range bounds {
        // I say c.Precision rather than bits because I think these should be equal
        // and if there is a bug, this will certainly break quicker.
        num.SetPrec(c.Precision)
    }
}

func (c *configurator) useNative() {
    c.NumericsStrategy = NativeNumericsMode
}

func (c *configurator) useBig() {
    c.NumericsStrategy = BigFloatNumericsMode
}

func (c *configurator) parseUserCoords() {
    bigActions := []func(*big.Float){
        func(realMin *big.Float) { c.RealMin = *realMin },
        func(realMax *big.Float) { c.RealMax = *realMax },
        func(imagMin *big.Float) { c.ImagMin = *imagMin },
        func(imagMax *big.Float) { c.ImagMax = *imagMax },
    }

    desc := c.UserRequest
    userInput := []string{
        desc.RealMin,
        desc.RealMax,
        desc.ImagMin,
        desc.ImagMax,
    }

    inputNames := []string{"realMin", "realMax", "imagMin", "imagMax"}

    for i, num := range userInput {
        bigFloat, bigErr := parseBig(num)

        if bigErr != nil {
            parsePanic(bigErr, inputNames[i])
        }

        // Handle parse results
        bigActions[i](bigFloat)
    }
}

// Choose an optimal strategy for rendering the image
func (c *configurator) chooseFastRenderStrategy() {
    req := c.UserRequest

    area := req.ImageWidth * req.ImageHeight
    numerics := c.NumericsStrategy

    if numerics == AutoDetectNumericsMode {
        log.Panic("Must choose render strategy after numerics system")
    }

    if area < DefaultTinyImageArea && numerics == NativeNumericsMode {
        // Use `SequenceRenderStrategy' when
        // We have native arithmetic and the image is tiny
        c.useSequenceRenderer()
    } else if req.Jobs <= DefaultLowThreading {
        // Use `RegionRenderStrategy' when
        // the number of jobs is small
        c.useRegionRenderer()
    } else {
        // Use `ConcurrentRegionRenderStrategy' otherwise
        c.useSharedRegionRenderer()
    }
}

func (c *configurator) useSequenceRenderer() {
    c.RenderStrategy = SequenceRenderMode
}

func (c *configurator) useRegionRenderer() {
    c.RenderStrategy = RegionRenderMode
}

func (c *configurator) useSharedRegionRenderer() {
    c.RenderStrategy = SharedRegionRenderMode
}

func (c *configurator) howManyBits() uint {
    // For now, always choose Native Arithmetic
    return 53
}