package region

import (
    "image"
    "functorama.com/demo/base"
    "functorama.com/demo/sequence"
    "functorama.com/demo/draw"
)

type RegionNumericsFactory interface {
    Build() RegionNumerics
}

// RegionNumerics provides rendering calculations for the "region" render strategy.
type RegionNumerics interface {
    base.OpaqueProxyFlyweight
    Rect() image.Rectangle
    EvaluateAllPoints(iterateLimit uint8)
    Split()
    OnGlitchCurve(iterateLimit uint8, glitchSamples uint) bool
    MandelbrotPoints() []base.MandelbrotMember
    RegionMember() base.MandelbrotMember
    Children() []RegionNumerics
    RegionSequence() ProxySequence
}

type ProxySequence interface {
    base.OpaqueProxyFlyweight
    sequence.SequenceNumerics
}

// RenderSequentialRegion takes a RegionNumerics but renders the region in a sequential
// (column-wise) manner
func RenderSequenceRegion(numerics RegionNumerics, context draw.DrawingContext, iterateLimit uint8) {
    smallNumerics := numerics.RegionSequence()
    smallNumerics.ClaimExtrinsics()
    smallNumerics.ImageDrawSequencer(context)
    smallNumerics.MandelbrotSequence(iterateLimit)
}

// SequenceCollapse is analogous to RenderSequentialRegion, but it returns the Mandelbrot render
// results rather than drawing them to the image.
func SequenceCollapse(numerics RegionNumerics, iterateLimit uint8) []base.PixelMember {
    sequence := numerics.RegionSequence()
    sequence.ClaimExtrinsics()
    sequence.MemberCaptureSequencer()
    sequence.MandelbrotSequence(iterateLimit)
    return sequence.CapturedMembers()
}

// Subdivide takes a RegionNumerics and tries to split the region into subregions.  It returns true
// if the subdivision occurred.  The subdivision won't occur if the region is Uniform or in an area
// where glitches are likely.
func Subdivide(numerics RegionNumerics, iterateLimit uint8, glitchSamples uint) bool {
    numerics.EvaluateAllPoints(iterateLimit)

    if !Uniform(numerics) || numerics.OnGlitchCurve(iterateLimit, glitchSamples) {
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

// Collapse returns true if the region is below the necessary size for subdivision
func Collapse(split RegionNumerics, collapseSize int) bool {
    rect := split.Rect()
    return rect.Dx() <= collapseSize || rect.Dy() <= collapseSize
}