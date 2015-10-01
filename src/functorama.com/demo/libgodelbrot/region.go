package libgodelbrot

import (
    "image"
)

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
        c: c,
    }
}

type Region struct {
    topLeft *EscapePoint
    topRight *EscapePoint
    bottomLeft *EscapePoint
    bottomRight *EscapePoint
    midPoint *EscapePoint
}

func NewRegion(topLeft complex128, bottomRight complex128) *Region {
    left := real(topLeft)
    right := real(bottomRight)
    top := imag(topLeft)
    bottom := imag(bottomRight)
    trPos := complex(right, top)
    blPos := complex(left, bottom)
    midPos := complex(
        left + ((right - left) / 2.0), 
        bottom + ((top - bottom) / 2.0),
    )

    tl := NewEscapePoint(topLeft)
    tr := NewEscapePoint(trPos)
    bl := NewEscapePoint(blPos)
    br := NewEscapePoint(bottomRight)
    mid := NewEscapePoint(midPos)

    return &Region{
        topLeft: tl,
        topRight: tr,
        bottomLeft: bl,
        bottomRight: br,
        midPoint: mid,
    }  
}

func (r Region) Points() []*EscapePoint {
    return []*EscapePoint{
        r.topLeft, 
        r.topRight, 
        r.bottomLeft, 
        r.bottomRight, 
        r.midPoint,
    }
}

type Subregion struct {
    populated bool
    children []*Region
}

func (r Region) Subdivide(iterateLimit uint8, divergeLimit float64) Subregion {
    points := r.Points()
    // Ensure points are all evaluated
    for _, p := range points {
        if (!p.evaluated) {
            p.membership =  Mandelbrot(p.c, iterateLimit, divergeLimit)
            p.evaluated = true
        }
    }

    // If inverse divergence on all points is the same, no need to subdivide
    divide := false
    last := points[0].membership.InvDivergence
    for _, p := range points[1:] {
        if p.membership.InvDivergence != last {
            divide = true
            break
        }
    }

    if divide {
        return r.Split();
    } else {
        return Subregion{
            populated: false,
        }
    }
}

func (r Region) Split() Subregion {
    topLeftPos := r.topLeft.c
    bottomRightPos := r.bottomRight.c
    midPos := r.midPoint.c

    left := real(topLeftPos)
    right := real(bottomRightPos)
    top := imag(topLeftPos)
    bottom := imag(bottomRightPos)
    midR := real(midPos)
    midI := imag(midPos)

    topSideMid := NewEscapePointReals(midR, top)
    bottomSideMid := NewEscapePointReals(midR, bottom)
    leftSideMid := NewEscapePointReals(left, midI)
    rightSideMid := NewEscapePointReals(right, midI)

    leftSectorMid := left + ((midR - left) / 2.0)
    rightSectorMid := midR + ((right - midR) / 2.0)
    topSectorMid := midI + ((top - midI) / 2.0)
    bottomSectorMid := bottom + ((midI - bottom) / 2.0)

    tl := Region{
        topLeft: r.topLeft,
        topRight: topSideMid,
        bottomLeft: leftSideMid,
        bottomRight: r.midPoint,
        midPoint: NewEscapePointReals(leftSectorMid, topSectorMid),
    }
    tr := Region{
        topLeft: topSideMid,
        topRight: r.topRight,
        bottomLeft: r.midPoint,
        bottomRight: rightSideMid,
        midPoint: NewEscapePointReals(rightSectorMid, topSectorMid),
    }
    bl := Region{
        topLeft: leftSideMid,
        topRight: r.midPoint,
        bottomLeft: r.bottomLeft,
        bottomRight: bottomSideMid,
        midPoint: NewEscapePointReals(leftSectorMid, bottomSectorMid),
    }
    br := Region{
        topLeft: r.midPoint,
        topRight: rightSideMid,
        bottomLeft: bottomSideMid,
        bottomRight: r.bottomRight,
        midPoint: NewEscapePointReals(rightSectorMid, bottomSectorMid),
    }

    return Subregion{
        populated: true,
        children: []*Region{&tl, &tr, &bl, &br},
    }
}

func (region Region) Rect(config *RenderConfig) image.Rectangle {
    l, t := config.PlaneToPixel(region.topLeft.c)
    r, b := config.PlaneToPixel(region.bottomRight.c)
    return image.Rect(int(l), int(t), int(r), int(b))
}

func (r Region) Collapse(config *RenderConfig) bool {
    rect := r.Rect(config)
    iCollapse := int(config.RegionCollapse)
    return rect.Dx() <= iCollapse || rect.Dy() <= iCollapse
}