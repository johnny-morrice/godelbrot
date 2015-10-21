package libgodelbrot

import (
	"image"
)

// This module covers internal details of the numeric strategies used by
// Godelbrot.
// The internal detail of each strategy is defined by an interfac.
// This is because the render control algorithms are also strategies that vary
// independently.

// NumericsSystem represents a full numerics stack that can interface with any Godelbrot render
// strategy
type NumericSystem interface {
	SharedRegionNumerics
	SharedSequentialNumerics
}

// SequentialNumerics provides sequential (column-wise) rendering calculations
type SequentialNumerics interface {
	OpaqueFlyweightProxy
	MandelbrotSequence(iterateLimit uint8)
	ImageDrawSequencer(draw DrawingContext)
	MemberCaptureSequencer()
	CapturedMembers() []PixelMember
}

// SharedSequentialNumerics provides sequential (column-wise) rendering calculations for a threaded
// render strategy
type SharedSequentialNumerics interface {
	SequentialNumerics
	OpaqueThreadPrototype
}

// RegionNumerics provides rendering calculations for the "region" render strategy.
type RegionNumerics interface {
	OpaqueProxyFlyweight
	Rect() image.Rectangle
	EvaluateAllPoints(iterateLimit int)
	Split()
	Uniform() bool
	OnGlitchCurve(iterateLimit uint8, glitchSamples uint) bool
	MandelbrotPoints() []MandelbrotMember
	RegionMember()
	Subdivide() bool
	Children() []RegionNumerics
	RegionalSequenceNumerics() SequentialNumerics
}

// SharedRegionNumerics provides a RegionNumerics for threaded render stregies
type SharedRegionNumerics interface {
	RegionNumerics
	OpaqueThreadPrototype
}

// DrawRegion combines a RegionNumerics and DrawingContext, implementing RegionDrawingContext
type DrawRegion struct {
	Numerics RegionNumerics
	Draw     DrawingContext
}

func (draw DrawRegion) RegionMember() MandelbrotMember {
	return draw.Numerics.RegionMember()
}

func (draw DrawRegion) Rect() image.Rectangle {
	return draw.Numerics.Rect()
}

func (draw DrawRegion) Picture() *image.NRGBA {
	return draw.Draw.Picture()
}

func (draw DrawRegion) Colors() Palette {
	return draw.Draw.Colors()
}

// PixelMember is a MandelbrotMember associated with a pixel
type PixelMember struct {
	I      int
	J      int
	member MandelbrotMember
}

// RenderSequentialRegion takes a RegionNumerics but renders the region in a sequential (column-wise) manner
//
func RenderSequentialRegion(numerics RegionNumerics) {
	smallNumerics := numerics.RegionalSequenceNumerics()
	smallNumerics.ClaimExtrinsics()
	SequentialRenderImage(smallNumerics)
}

// Subdivide takes a RegionNumerics and tries to split the region into subregions.  It returns true
// if the subdivision occurred.  The subdivision won't occur if the region is Uniform or in an area
// where glitches are likely.
func Subdivide(numerics RegionNumerics) bool {
	numerics.EvaluateAllPoints()

	if !Uniform(numerics) || numerics.OnGlitchCurve() {
		numerics.Split()
		return true
	}
	return false
}

// Uniform returns true if the region has the same Mandelbrot escape value across its bounds.
// Note this function performs no evaluation of membership.
func Uniform(numerics RegionNumerics) bool {
	// If inverse divergence on all points is the same, no need to subdivide
	points := numerics.MandelbrotPoints()
	first := points[0].InverseDivergence()
	for _, p := range points[1:] {
		if p.InverseDivergence() != first {
			return false
		}
	}
	return true
}

// SequenceCollapse is analogous to RenderSequence region, but it returns the Mandelbrot render
// results rather than drawing them to the image.
func SequenceCollapse(numerics RegionNumerics) []PixelMember {
	sequence := region.RegionalSequenceNumerics()
	sequence.ClaimExtrinsics()
	sequence.MemberCaptureSequencer()
	MandelbrotSequence(smallConfig)
	return sequence.CapturedMembers()
}

// Collapse returns true if the region is below the necessary size for subdivision
func Collapse(split RegionNumerics, collapseSize int) bool {
	rect := split.Rect()
	return rect.Dx() <= collapseSize || rect.Dy() <= collapseSize
}
