package libgodelbrot

import (
    "math/big"
    "log"
    "fmt"
)

// Object to initialize the godelbrot system
type configurator Info

// InitializeContext examines the description, chooses a renderer, numerical system and palette.
func Configure(req *Request) (*Info, error) {
    c := &configurator{}
    c.UserRequest = *req

    nerr := c.chooseNumerics()

    if nerr != nil {
        return nil, nerr
    }

    rerr := c.chooseRenderStrategy()

    if rerr != nil {
        return nil, rerr
    }

    perr := c.choosePalette()

    if perr != nil {
        return nil, perr
    }

    return (*Info)(c), nil
}

// Initialize the render system
func (c *configurator) chooseRenderStrategy() error {
    req := c.UserRequest
    switch req.Renderer {
    case AutoDetectRenderMode:
        c.chooseFastRenderStrategy()
    case SequenceRenderMode:
        c.useSequenceRenderer()
    case RegionRenderMode:
        c.useRegionRenderer()
    default:
        return fmt.Errorf("Unknown render mode: %v", req.Renderer)
    }

    return nil
}

// Initialize the numerics system
func (c *configurator) chooseNumerics() error {
    desc := c.UserRequest
    perr := c.parseUserCoords()

    if perr != nil {
        return perr
    }

    switch desc.Numerics {
    case AutoDetectNumericsMode:
        c.chooseAccurateNumerics()
    case NativeNumericsMode:
        c.useNative()
        c.Precision = 53
        c.usePrec()
    case BigFloatNumericsMode:
        c.selectUserPrec()
        c.usePrec()
        c.useBig()
    default:
        return fmt.Errorf("Unknown numerics mode:", desc.Numerics)
    }

    return nil
}

func (c *configurator) selectUserPrec() {
    userPrec := c.UserRequest.Precision
    if userPrec > 0 {
        c.Precision = userPrec
    } else {
        c.Precision = c.howManyBits()
    }
}

func (c *configurator) chooseAccurateNumerics() {
    // 53 bits precision is available to 64 bit floats
    const prec64 uint = 53

    c.selectUserPrec()
    c.usePrec()
    if c.Precision > prec64 {
        c.useBig()
    } else {
        c.useNative()
    }
}

func (c *configurator) usePrec() {
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

func (c *configurator) parseUserCoords() error {
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
            return fmt.Errorf("Could not parse %v: %v", inputNames[i], bigErr)
        }

        // Handle parse results
        bigActions[i](bigFloat)
    }

    return nil
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
    } else {
        c.useRegionRenderer()
    }
}

func (c *configurator) useSequenceRenderer() {
    c.RenderStrategy = SequenceRenderMode
}

func (c *configurator) useRegionRenderer() {
    c.RenderStrategy = RegionRenderMode
}

func (c *configurator) howManyBits() uint {
    // For now, always choose Native Arithmetic
    return 53
}

func (c *configurator) choosePalette() error {
    code := c.UserRequest.PaletteCode
    switch code {
    case "pretty":
        c.PaletteType = Pretty
    case "redscale":
        c.PaletteType = Redscale
    case "grayscale":
        c.PaletteType = Grayscale
    default:
        return fmt.Errorf("Invalid palette code: %v", code)
    }

    return nil
}