package config

import (
	"errors"
)

type AspectConservation uint8

const (
	Stretch = AspectConservation(iota)
	Shrink
	Grow
)

// Request is a user description of the render to be accomplished
type Request struct {
	IterateLimit uint8
	DivergeLimit float64
	RealMin      string
	RealMax      string
	ImagMin      string
	ImagMax      string
	ImageWidth   uint
	ImageHeight  uint
	PaletteCode  string
	FixAspect    AspectConservation
	// Render algorithm
	Renderer RenderMode
	// Number of render threads
	Jobs           uint16
	RegionCollapse uint
	// Numerical system
	Numerics NumericsMode
	// Number of samples taken when detecting region render glitches
	RegionSamples uint
	// Number of bits for big.Float rendering
	Precision uint
}

// Available render algorithms
type RenderMode uint

const (
	AutoDetectRenderMode = RenderMode(iota)
	RegionRenderMode
	SequenceRenderMode
)

// Available numeric systems
type NumericsMode uint

const (
	// Functions should auto-detect the correct system for rendering
	AutoDetectNumericsMode = NumericsMode(iota)
	// Use the native CPU arithmetic operations
	NativeNumericsMode
	// Use arithmetic based around the standard library big.Float type
	BigFloatNumericsMode
)

type ZoomBounds struct {
	Xmin uint
	Xmax uint
	Ymin uint
	Ymax uint
}

func (zb *ZoomBounds) Validate() error {
	if zb.Xmin >= zb.Xmax || zb.Ymin >= zb.Ymax {
		return errors.New("Min and max zoom boundaries are invalid.")
	}
	return nil
}

type ZoomTarget struct {
	ZoomBounds
	// Reconsider numerical system and render modes as appropriate.
	Reconfigure bool
	// Increase precision.  With Reconfigure, this should automatically engage arbitrary
	// precision mode.
	UpPrec bool
	// Number of frames for zoom
	Frames uint
}
