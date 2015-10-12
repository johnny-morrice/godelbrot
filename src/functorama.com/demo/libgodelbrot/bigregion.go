package libgodelbrot

// This module sucks because it is a near total clone of the native version

import (
    "image"
    "math/big"
)

type bigSubregion struct {
    populated bool
    children  []BigRegion
}

type BigRegionNumerics struct {
    Collapser
    BigBaseNumerics
    region bigRegion
    subregion bigSubregion
    heap *BigMandelbrotThunkHeap
}

type BigRegion struct {
    topLeft     *BigMandelbrotThunk
    topRight    *BigMandelbrotThunk
    bottomLeft  *BigMandelbrotThunk
    bottomRight *BigMandelbrotThunk
    midPoint    *BigMandelbrotThunk
}

func (bigFloat *BigRegionNumerics) MandelbrotPoints() {
    r := bigFloat.Region
    return []MandelbrotMember {
        r.topLeft.membership,
        r.topRight.membership,
        r.bottomLeft.membership,
        r.bottomRight.membership,
        r.midPoint.membership,
    }
}

func (bigFloat *BigRegionNumerics) EvaluateAllPoints() {
    points := []*BigMandelbrotThunk{
        r.topLeft,
        r.topRight,
        r.bottomLeft,
        r.bottomRight,
        r.midPoint,
    }
    // Ensure points are all evaluated
    for _, p := range points {
        if !p.evaluated {
            EvalThunk(p)
        }
    }
}

// A glitch is possible when points are uniform near the set
// Due to the shape of the set, a rectangular Bigregion is not a good approximation
// An anologous glitch happens when the entire Bigregion is much larger than the set
// We handle both these cases here
func (bigFloat *BigRegionNumerics) OnGlitchCurve() bool {
    r := bigFloat.Region
    member := bigFloat.RegionMember()
    iDiv := member.InvDivergence()
    iLimit, dLimit := bigFloat.MandelbrotLimits()
    if iDiv == 0 || iDiv == 1 || member.InSet() {
        sqrtChecks := bigFloat.GlitchSamples()
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
                &checkMember.Mandelbrot(iLimit, dLimit)
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

func (bigFloat *BigRegionNumerics) Split() {
    heap := bigFloat.Heap
    r := bigFloat.Region

    topLeftPos := r.topLeft.c
    bottomRightPos := r.bottomRight.c
    midPos := r.midPoint.c

    left := topLeftPos.Real()
    right := bottomRightPos.Real()
    top := topLeftPos.Imag()
    bottom := bottomRightPos.Imag()
    midR := midPos.Real()
    midI := midPos.Imag()

    topSideMid := heap.BigMandelbrotThunk(midR, top)
    bottomSideMid := heap.BigMandelbrotThunk(midR, bottom)
    leftSideMid := heap.BigMandelbrotThunk(left, midI)
    rightSideMid := heap.BigMandelbrotThunk(right, midI)

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
        midPoint:    heap.BigMandelbrotThunk(leftSectorMid, topSectorMid),
    }
    tr := BigRegion{
        topLeft:     topSideMid,
        topRight:    r.topRight,
        bottomLeft:  r.midPoint,
        bottomRight: rightSideMid,
        midPoint:    heap.BigMandelbrotThunk(rightSectorMid, topSectorMid),
    }
    bl := BigRegion{
        topLeft:     leftSideMid,
        topRight:    r.midPoint,
        bottomLeft:  r.bottomLeft,
        bottomRight: bottomSideMid,
        midPoint:    heap.BigMandelbrotThunk(leftSectorMid, bottomSectorMid),
    }
    br := BigRegion{
        topLeft:     r.midPoint,
        topRight:    rightSideMid,
        bottomLeft:  bottomSideMid,
        bottomRight: r.bottomRight,
        midPoint:    heap.BigMandelbrotThunk(rightSectorMid, bottomSectorMid),
    }

    bigFloat.Subregion = BigSubregion{
        populated: true,
        children:  []BigRegion{tl, tr, bl, br},
    }
}

func (bigFloat *BigRegionNumerics) Rect() image.Rectangle {
    l, t := bigFloat.PlaneToPixel(Bigregion.topLeft.c)
    r, b := bigFloat.PlaneToPixel(Bigregion.bottomRight.c)
    return image.Rect(int(l), int(t), int(r), int(b))
}

// Return MandelbrotMember
// Does not check if the region's thunks have been evaluated
func (bigFloat *BigRegionNumerics) RegionMember() MandelbrotMember {
    return bigFloat.Region.topLeft.member.MandelbrotMember
}
