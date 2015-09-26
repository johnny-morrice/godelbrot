package libgodelbrot

import {
    "math/cmplx"
}

type EscapePoint struct {
    evaluated bool
    c complex128
    membership MandelbrotMember
}

func NewEscapePointReals(r float64, i float64) *EscapePoint {
    return NewEscapePoint(complex(r, i))
}

func NewEscapePoint(c complex128) *EscapePoint {
    return &EscapePoint{
        evaluated: false,
        c: c
    }
}

type Region struct {
    topLeft *EscapePoint
    topRight *EscapePoint
    bottomLeft *EscapePoint
    bottomRight *EscapePoint
    midPoint *EscapePoint
}

type Subregion struct {
    children []Region
}

func (r *Region) Subdivide(iterateLimit uint8, divergeLimit uint8) Subregion {
    points := []{r.topLeft, r.topRight, r.bottomLeft, r.bottomRight, r.midPoint}
    // Ensure points are all evaluated
    for p := range points {
        if (!p.evaluated) {
            p.membership =  Mandelbrot(p.c, iterateLimit, divergeLimit)
            p.evaluated = true
        }
    }

    // If inverse divergence on all points is the same, no need to subdivide
    divide := false
    last := points[0].membership.InvDivergence
    for p := range points[1:] {
        if p.membership.InvDivergence != last {
            divide = true
            break
        }
    }

    if divide {
        return r.Split();
    } else {
        return nil
    }
}

func (r Region) Split() Subregion {
    topLeftPos := r.topLeft.c
    bottomRightPos := r.bottomRight.c
    midPos := r.midPoint.c

    left := cmplx.Real(topLeftPos)
    right := cmplx.Real(bottomRightPos)
    top := cmplx.Imag(topLeftPos)
    bottom := cmplx.Imag(bottomRightPos)
    midR := cmplx.Real(midPos)
    midI := cmplx.Imag(midPos)

    topSideMid := NewEscapePointReals(midR, top)
    bottomSideMid := NewEscapePointReals(midR, bottom)
    leftSideMid := NewEscapePointReals(left, midI)
    rightSideMid := NewEscapePointReals(right, midI)

    leftSectorMid := (midR - left) / 2.0
    rightSectorMid := (right - midR) / 2.0
    topSectorMid := (midI - top) / 2.0
    bottomSectorMid := (bottom - midI) / 2.0

    tl := Region{
        topLeft: r.topLeft,
        topRight: topSideMid,
        bottomLeft: leftSideMid,
        bottomRight: r.midPoint,
        midPoint: NewEscapePointReals(leftSectorMid, topSectorMid)
    }
    tr := Region{
        topLeft: topSideMid,
        topRight: r.topRight,
        bottomLeft: r.midPoint,
        bottomRight: rightSideMid,
        midPoint: NewEscapePointReals(rightSectorMid, topSectorMid)
    }
    bl := Region{
        topLeft: leftSideMid,
        topRight: r.midPoint,
        bottomLeft: r.bottomLeft,
        bottomRight: bottomSideMid,
        midPoint: NewEscapePointReals(leftSectorMid, bottomSectorMid)
    }
    br := Region{
        topLeft: r.midPoint,
        topRight: rightSideMid,
        bottomLeft: bottomSideMid,
        bottomRight: r.bottomRight,
        midPoint: NewEscapePointReals(rightSectorMid, bottomSectorMid)
    }

    return []{tl, tr, bl, br}
}