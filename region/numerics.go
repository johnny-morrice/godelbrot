package region

import (
    "image"
    "github.com/johnny-morrice/godelbrot/base"
    "github.com/johnny-morrice/godelbrot/sequence"
    "github.com/johnny-morrice/godelbrot/draw"
)

type Subdivider interface {
    Rect() image.Rectangle
    Split()
    MandelbrotPoints() []base.EscapeValue
    SampleDivs() (<-chan uint8, chan<- bool)
    RegionMember() base.EscapeValue
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
func RenderSequenceRegion(reg RegionNumerics, ctx draw.DrawingContext) {
    reg.ClaimExtrinsics()
    seq := reg.RegionSequence()
    seq.Extrinsically(func () {
        sequence.ImageSequence(seq, ctx)
    })
}

// SequenceCollapse is analogous to RenderSequentialRegion, but it returns the Mandelbrot render
// results rather than drawing them to the image.
func SequenceCollapse(num RegionNumerics) []base.PixelMember {
    seq := num.RegionSequence()
    var pix []base.PixelMember
    seq.Extrinsically(func () {
        pix = sequence.Capture(seq)
    })
    return pix
}

// Subdivide takes a RegionNumerics and tries to split the region into subregions.  It returns true
// if the subdivision occurred.  The subdivision won't occur if the region is Uniform or in an area
// where glitches are likely.
func Subdivide(reg RegionNumerics) bool {
    if !Uniform(reg) {
        reg.Split()
        return true
    }
    return false
}

// Uniform returns true if the region has the same Mandelbrot escape value across its bounds.
func Uniform(reg RegionNumerics) bool {
    // If inverse divergence on all points is the same, no need to subdivide
    idivs, done := reg.SampleDivs()
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
func Collapse(reg RegionNumerics, sizelim int) bool {
    rect := reg.Rect()
    return rect.Dx() <= sizelim || rect.Dy() <= sizelim
}