package godelbrot

import (
	"fmt"
	"github.com/johnny-morrice/godelbrot/config"
	bb "github.com/johnny-morrice/godelbrot/internal/bigbase"
	"log"
	"math/big"
)

// Object to initialize the godelbrot system
type configurator Info

// InitializeContext examines the description, chooses a renderer, numerical system and palette.
func Configure(req *config.Request) (*Info, error) {
	c := &configurator{}
	c.UserRequest = *req

	nerr := c.chooseNumerics()

	if nerr != nil {
		return nil, nerr
	}

	verr := c.validate()

	if verr != nil {
		return nil, verr
	}

	rerr := c.chooseRenderStrategy()

	if rerr != nil {
		return nil, rerr
	}

	perr := c.choosePalette()

	if perr != nil {
		return nil, perr
	}

	if req.FixAspect != config.Stretch {
		ferr := c.fixAspect()
		if ferr != nil {
			return nil, ferr
		}
	}

	return (*Info)(c), nil
}

func (c *configurator) validate() error {
	rcmp := c.RealMin.Cmp(&c.RealMax)
	icmp := c.ImagMin.Cmp(&c.ImagMax)

	if rcmp == 0 || rcmp == 1 || icmp == 0 || icmp == 1 {
		return fmt.Errorf("Invalid bounds")
	}

	return nil
}

func (c *configurator) fixAspect() error {
	if __DEBUG {
		log.Println("Fixing aspect ratio")
	}

	rmin := big.NewFloat(0.0).Copy(&c.RealMin)
	rmax := big.NewFloat(0.0).Copy(&c.RealMax)
	imin := big.NewFloat(0.0).Copy(&c.ImagMin)
	imax := big.NewFloat(0.0).Copy(&c.ImagMax)

	planeWidth := bb.MakeBigFloat(0.0, c.Precision)
	planeHeight := bb.MakeBigFloat(0.0, c.Precision)
	planeAspect := bb.MakeBigFloat(0.0, c.Precision)

	planeWidth.Sub(rmax, rmin)
	planeHeight.Sub(imax, imin)
	planeAspect.Quo(&planeWidth, &planeHeight)

	nativePictureAspect := float64(c.UserRequest.ImageWidth) / float64(c.UserRequest.ImageHeight)
	pictureAspect := bb.MakeBigFloat(nativePictureAspect, c.Precision)
	thindicator := planeAspect.Cmp(&pictureAspect)

	adjustWidth := func() {
		trans := bb.MakeBigFloat(0.0, c.Precision)
		trans.Mul(&planeHeight, &pictureAspect)
		rmax.Add(rmin, &trans)
		c.RealMax = *rmax
	}

	adjustHeight := func() {
		trans := bb.MakeBigFloat(0.0, c.Precision)
		trans.Quo(&planeWidth, &pictureAspect)
		imax.Add(imin, &trans)
		c.ImagMax = *imax
	}

	if __TRACE {
		log.Printf("Old Plane Width: %v", bb.DbgF(planeWidth))
		log.Printf("Old Plane Height: %v", bb.DbgF(planeHeight))
		log.Printf("Old Plane Aspect: %v", bb.DbgF(planeAspect))
		log.Printf("Image Width: %v", c.UserRequest.ImageWidth)
		log.Printf("Image Height: %v", c.UserRequest.ImageHeight)
		log.Printf("Image Aspect: %v", bb.DbgF(pictureAspect))
	}

	var strategy = c.UserRequest.FixAspect

	if strategy == config.Grow {
		// Then the plane is too short, so must be made taller
		if thindicator == 1 {
			if __DEBUG {
				log.Println("Plane becomes taller")
			}
			adjustHeight()
		} else if thindicator == -1 {
			if __DEBUG {
				log.Println("Plane becomes fatter")
			}
			// Then the plane is too thin, and must be made fatter
			adjustWidth()
		}
	} else if strategy == config.Shrink {
		if thindicator == 1 {
			// The plane is too fat, so must be made thinner
			if __DEBUG {
				log.Println("Plane becomes thinner")
			}
			adjustWidth()
		} else if thindicator == -1 {
			if __DEBUG {
				log.Println("Plane becomes shorter")
			}
			// The plane is too tall, so must be made shorter
			adjustHeight()
		}
	}

	if __TRACE && (strategy == config.Shrink || strategy == config.Grow) {
		fixWidth := bb.MakeBigFloat(0.0, c.Precision)
		fixHeight := bb.MakeBigFloat(0.0, c.Precision)
		fixAspect := bb.MakeBigFloat(0.0, c.Precision)

		fixWidth.Sub(rmax, rmin)
		fixHeight.Sub(imax, imin)
		fixAspect.Quo(&fixWidth, &fixHeight)

		log.Printf("New Plane Width: %v", bb.DbgF(fixWidth))
		log.Printf("New Plane Height: %v", bb.DbgF(fixHeight))
		log.Printf("Fixed Aspect: %v", bb.DbgF(fixAspect))
	}

	return nil
}

// Initialize the render system
func (c *configurator) chooseRenderStrategy() error {
	req := c.UserRequest
	switch req.Renderer {
	case config.AutoDetectRenderMode:
		c.chooseFastRenderStrategy()
	case config.SequenceRenderMode:
		c.useSequenceRenderer()
	case config.RegionRenderMode:
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
	case config.AutoDetectNumericsMode:
		c.chooseAccurateNumerics()
	case config.NativeNumericsMode:
		c.useNative()
		c.Precision = 53
		c.usePrec()
	case config.BigFloatNumericsMode:
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
	c.NumericsStrategy = config.NativeNumericsMode
}

func (c *configurator) useBig() {
	c.NumericsStrategy = config.BigFloatNumericsMode
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

	if numerics == config.AutoDetectNumericsMode {
		log.Panic("Must choose render strategy after numerics system")
	}

	bigsz := area > DefaultTinyImageArea
	weirdbase := numerics != config.NativeNumericsMode
	squarepic := req.ImageWidth == req.ImageHeight

	if (bigsz || weirdbase) && squarepic {
		c.useRegionRenderer()
	} else {
		c.useRegionRenderer()
	}
}

func (c *configurator) useSequenceRenderer() {
	c.RenderStrategy = config.SequenceRenderMode
}

func (c *configurator) useRegionRenderer() {
	c.RenderStrategy = config.RegionRenderMode
}

// Sample method to discover how many bits needed
func (c *configurator) howManyBits() uint {
	bounds := []big.Float{
		c.RealMin,
		c.RealMax,
		c.ImagMin,
		c.ImagMax,
	}

	bits := uint(0)
	for _, bnd := range bounds {
		prec := bnd.MinPrec()
		if prec > bits {
			bits = prec
		}
	}

	return bits
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
