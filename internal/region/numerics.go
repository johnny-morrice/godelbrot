package region

import (
    "image"
    "github.com/johnny-morrice/godelbrot/internal/base"
    "github.com/johnny-morrice/godelbrot/internal/sequence"
    "github.com/johnny-morrice/godelbrot/internal/draw"
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

type Region struct {
    Xmin int
    Xmax int
    Ymin int
    Ymax int
}

func InitRegion(bn *base.BaseNumerics) Region {
    r := Region{}
    r.Xmin = bn.PicXMin
    r.Xmax = bn.PicXMax
    r.Ymin = bn.PicYMin
    r.Ymax = bn.PicYMax

    return r
}

func (r Region) Split() []Region {
    width := r.Xmax - r.Xmin
    height := r.Ymax - r.Ymin

    toxmid := r.Xmin + (width / 2)
    toymid := r.Ymin + (height / 2)

    fromxmid := toxmid + 1
    fromymid := toymid + 1 

    tl := Region{}
    tl.Xmin = r.Xmin
    tl.Ymin = r.Ymin
    tl.Xmax = toxmid
    tl.Ymax = toymid

    tr := Region{}
    tr.Xmin = fromxmid
    tr.Xmax = r.Xmax
    tr.Ymin = r.Ymin
    tr.Ymax = toymid

    bl := Region{}
    bl.Xmin = r.Xmin
    bl.Xmax = toxmid
    bl.Ymin = fromymid
    bl.Ymax = r.Ymax

    br := Region{}
    br.Xmin = fromxmid
    br.Xmax = r.Xmax
    br.Ymin = fromymid
    br.Ymax = r.Ymax

    return []Region{tl, tr, bl, br}
}