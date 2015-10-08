package libgodelbrot

// This module sucks because it is a near total clone of the native version

import (
    "image"
    "math/big"
)

type BigSubregion struct {
    populated bool
    children  []BigRegion
}

type BigRegionRenderContext struct {
    Region BigRegion
    Subregion BigSubregion
    Config *BigConfig
    Heap *BigEscapePointHeap
}

type BigRegion struct {
    topLeft     *BigEscapePoint
    topRight    *BigEscapePoint
    bottomLeft  *BigEscapePoint
    bottomRight *BigEscapePoint
    midPoint    *BigEscapePoint
}

func CreateBigRegion(topLeft BigComplex, bottomRight BigComplex) BigRegion {
    left := topLeft.Real()
    right := bottomRight.Real()
    top := topLeft.Imag()
    bottom := bottomRight.Imag()
    trPos := BigComplex{R: right, I: top}
    blPos := BigComplex{R: left, I: bottom}

    midPos := BigComplex{}
    
    midPos.R = right.Copy()
    midPos.R.Add(midPos.R, left)
    midPos.R.Quo(midPos.R, bigTwo)

    midPos.I = top.Copy()
    midPos.I.Add(midPos.I, bottom)
    midPos.I.Quo(midPos.I, bigTwo)

    tl := NewBigEscapePoint(topLeft)
    tr := NewBigEscapePoint(trPos)
    bl := NewBigEscapePoint(blPos)
    br := NewBigEscapePoint(bottomRight)
    mid := NewBigEscapePoint(midPos)

    return BigRegion{
        topLeft:     tl,
        topRight:    tr,
        bottomLeft:  bl,
        bottomRight: br,
        midPoint:    mid,
    }
}

func WholeBigRegion(config *BigConfig) BigRegion {
    return CreateBigRegion(config.PlaneTopLeft(), config.PlaneBottomRight())
}

func (context *BigRegionRenderContext) MandelbrotPoints() {
    r := context.Region
    return []MandelbrotMember {
        r.topLeft.membership,
        r.topRight.membership,
        r.bottomLeft.membership,
        r.bottomRight.membership,
        r.midPoint.membership,
    }
}

func (context *BigRegionRenderContext) EvaluateAllPoints() {
    points := []*BigEscapePoint{
        r.topLeft,
        r.topRight,
        r.bottomLeft,
        r.bottomRight,
        r.midPoint,
    }
    // Ensure points are all evaluated
    for _, p := range points {
        if !p.evaluated {
            p.membership.C = p.c
            (&p.membership).Mandelbrot(config.IterateLimit, config.DivergeLimit)
            p.evaluated = true
        }
    }
}

// A glitch is possible when points are uniform near the set
// Due to the shape of the set, a rectangular Bigregion is not a good approximation
// An anologous glitch happens when the entire Bigregion is much larger than the set
// We handle both these cases here
func (context *BigRegionRenderContext) OnGlitchCurve() bool {
    r := context.Region
    config := context.Config
    member := r.topLeft.membership
    iDiv := member.InvDivergence
    if iDiv == 0 || iDiv == 1 || member.InSet {
        sqrtChecks := context.Config.GlitchSamples
        tl := r.topLeft.c
        br := r.bottomRight.c

        hUnit := br.Real().Copy()
        hUnit.Sub(hUnit, tl.Real())
        hUnit.Quo(hUnit, sqrtChecks)
        vUnit := tl.Imag().Copy()
        vUnit.Sub(h, br.Imag())
        vUnit.Quo(vUnit, sqrtChecks)

        x := tl.Real()
        for i := 0; i < sqrtChecks; i++ {
            y := tl.Imag()
            for j := 0; j < sqrtChecks; j++ {
                checkMember := BigMandelbrotMember {
                    C: BigComplex{R: x, I: y},
                }
                &checkMember.Mandelbrot(config.IterateLimit, config.DivergeLimit)
                if member.InvDivergence != iDiv {
                    return true
                }
                y.Sub(y, vUnit)
            }
            x.Add(x, hUnit)
        }
    }

    return false
}

func (context *BigRegionRenderContext) Split() {
    heap := context.Heap
    r := context.Region

    topLeftPos := r.topLeft.c
    bottomRightPos := r.bottomRight.c
    midPos := r.midPoint.c

    left := topLeftPos.Real()
    right := bottomRightPos.Real()
    top := topLeftPos.Imag()
    bottom := bottomRightPos.Imag()
    midR := midPos.Real()
    midI := midPos.Imag()

    topSideMid := heap.BigEscapePoint(midR, top)
    bottomSideMid := heap.BigEscapePoint(midR, bottom)
    leftSideMid := heap.BigEscapePoint(left, midI)
    rightSideMid := heap.BigEscapePoint(right, midI)

    leftSectorMid := left.Copy()
    leftSectorMid.Add(leftSectorMid, midR)
    leftSectorMid.Quot(leftSectorMid, bigTwo)

    rightSectorMid := right.Copy()
    rightSectorMid.Add(rightSectorMid, midR)
    rightSectorMid.Quo(rightSectorMid, bigTwo)

    topSectorMid := top.Copy() 
    topSectorMid.Add(topSectorMid, midI)
    topSectorMid.Quo(topSectorMid, bigTwo)

    bottomSectorMid := bottom.Copy()
    bottomSectorMid.Add(bottomSectorMid, midI)
    bottomSectorMid.Quo(bottomSectorMid, bigTwo)

    tl := BigRegion{
        topLeft:     r.topLeft,
        topRight:    topSideMid,
        bottomLeft:  leftSideMid,
        bottomRight: r.midPoint,
        midPoint:    heap.BigEscapePoint(leftSectorMid, topSectorMid),
    }
    tr := BigRegion{
        topLeft:     topSideMid,
        topRight:    r.topRight,
        bottomLeft:  r.midPoint,
        bottomRight: rightSideMid,
        midPoint:    heap.BigEscapePoint(rightSectorMid, topSectorMid),
    }
    bl := BigRegion{
        topLeft:     leftSideMid,
        topRight:    r.midPoint,
        bottomLeft:  r.bottomLeft,
        bottomRight: bottomSideMid,
        midPoint:    heap.BigEscapePoint(leftSectorMid, bottomSectorMid),
    }
    br := BigRegion{
        topLeft:     r.midPoint,
        topRight:    rightSideMid,
        bottomLeft:  bottomSideMid,
        bottomRight: r.bottomRight,
        midPoint:    heap.BigEscapePoint(rightSectorMid, bottomSectorMid),
    }

    context.Subregion = BigSubregion{
        populated: true,
        children:  []BigRegion{tl, tr, bl, br},
    }
}

func (context *BigRegionRenderContext) Rect() image.Rectangle {
    l, t := context.Config.PlaneToPixel(Bigregion.topLeft.c)
    r, b := context.Config.PlaneToPixel(Bigregion.bottomRight.c)
    return image.Rect(int(l), int(t), int(r), int(b))
}

func (context *BigRegionRenderContext) Collapse() bool {
    rect := context.rect
    iCollapse := int(context.Config.BigRegionCollapse)
    return rect.Dx() <= iCollapse || rect.Dy() <= iCollapse
}
