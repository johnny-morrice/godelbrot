package region

import (
    "image"
    "functorama.com/demo/base"
    "functorama.com/demo/sequence"
    "functorama.com/demo/draw"
)

type Subdivider interface {
    Rect() image.Rectangle
    Split()
    MandelbrotPoints() []base.MandelbrotMember
    SampleDivs() (<-chan uint8, chan<- bool)
    RegionMember() base.MandelbrotMember
    RegionSequence() ProxySequence
}

// RegionNumerics provides rendering calculations for the "region" render strategy.
type RegionNumerics interface {
    base.OpaqueProxyFlyweight
    Subdivider
    Children() []RegionNumerics // Bad to include method in extended interface, but how else?
}

type ProxySequence interface {
    base.OpaqueProxyFlyweight
    sequence.SequenceNumerics
}

// RenderSequentialRegion takes a RegionNumerics but renders the region in a sequential
// (column-wise) manner
func RenderSequenceRegion(numerics RegionNumerics, context draw.DrawingContext) {
    numerics.ClaimExtrinsics()
    smallNumerics := numerics.RegionSequence()
    smallNumerics.Extrinsically(func () {
        sequence.ImageSequence(smallNumerics, context)
    })
}

// SequenceCollapse is analogous to RenderSequentialRegion, but it returns the Mandelbrot render
// results rather than drawing them to the image.
func SequenceCollapse(numerics RegionNumerics) []base.PixelMember {
    collapse := numerics.RegionSequence()
    var seq []base.PixelMember
    collapse.Extrinsically(func () {
        seq = sequence.Capture(collapse)
    })
    return seq
}

// Subdivide takes a RegionNumerics and tries to split the region into subregions.  It returns true
// if the subdivision occurred.  The subdivision won't occur if the region is Uniform or in an area
// where glitches are likely.
func Subdivide(numerics RegionNumerics) bool {
    if !Uniform(numerics) {
        numerics.Split()
        return true
    }
    return false
}

// Uniform returns true if the region has the same Mandelbrot escape value across its bounds.
func Uniform(numerics RegionNumerics) bool {
    // If inverse divergence on all points is the same, no need to subdivide
    idivs, done := numerics.SampleDivs()
    first := <-idivs
    uni := true
    for d := range idivs {
        if d != first && uni {
            done<- true
            uni = false
        }
    }
    return uni
}

// Collapse returns true if the region is below the necessary size for subdivision
func Collapse(split RegionNumerics, collapseSize int) bool {
    rect := split.Rect()
    return rect.Dx() <= collapseSize || rect.Dy() <= collapseSize
}