package libgodelbrot

import (
	"image"
)

type NativeSubregion struct {
	populated bool
	children  []NativeRegion
}

type NativeRegionRenderContext struct {
	Region NativeRegion
	Subregion NativeSubregion
	Config *NativeConfig
	Heap *NativeEscapePointHeap
}

type NativeRegion struct {
	topLeft     *NativeEscapePoint
	topRight    *NativeEscapePoint
	bottomLeft  *NativeEscapePoint
	bottomRight *NativeEscapePoint
	midPoint    *NativeEscapePoint
}

func CreateNativeRegion(topLeft complex128, bottomRight complex128) NativeRegion {
	left := real(topLeft)
	right := real(bottomRight)
	top := imag(topLeft)
	bottom := imag(bottomRight)
	trPos := complex(right, top)
	blPos := complex(left, bottom)
	midPos := complex(
		(right+left) / 2.0,
		(top+bottom) / 2.0,
	)

	tl := NewNativeEscapePoint(topLeft)
	tr := NewNativeEscapePoint(trPos)
	bl := NewNativeEscapePoint(blPos)
	br := NewNativeEscapePoint(bottomRight)
	mid := NewNativeEscapePoint(midPos)

	return NativeRegion{
		topLeft:     tl,
		topRight:    tr,
		bottomLeft:  bl,
		bottomRight: br,
		midPoint:    mid,
	}
}

func WholeNativeRegion(config *NativeConfig) NativeRegion {
	return CreateNativeRegion(config.PlaneTopLeft(), config.PlaneBottomRight())
}

func (context *NativeRegionRenderContext) MandelbrotPoints() {
	r := context.Region
	return []MandelbrotMember {
		r.topLeft.membership,
		r.topRight.membership,
		r.bottomLeft.membership,
		r.bottomRight.membership,
		r.midPoint.membership,
	}
}

func (context *NativeRegionRenderContext) EvaluateAllPoints() {
	r := context.Region
    points := []*NativeEscapePoint{
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
// Due to the shape of the set, a rectangular Nativeregion is not a good approximation
// An anologous glitch happens when the entire Nativeregion is much larger than the set
// We handle both these cases here
func (context *NativeRegionRenderContext) OnGlitchCurve() bool {
	r := context.Region
	config := context.Config
	member := r.topLeft.membership
	iDiv := member.InvDivergence
	if iDiv == 0 || iDiv == 1 || member.InSet {
		sqrtChecks := 10
		sqrtChecksF := float64(sqrtChecks)
		tl := r.topLeft.c
		br := r.bottomRight.c
		w := real(br) - real(tl)
		h := imag(tl) - imag(br)
		vUnit := h / sqrtChecksF
		hUnit := w / sqrtChecksF
		x := real(tl)
		for i := 0; i < sqrtChecks; i++ {
			y := imag(tl)
			for j := 0; j < sqrtChecks; j++ {
				checkMember := NativeMandelbrotMember {
					C: complex(x, y),
				}
				&checkMember.Mandelbrot(config.IterateLimit, config.DivergeLimit)
				if member.InvDivergence != iDiv {
					return true
				}
				y -= vUnit
			}
			x += hUnit
		}
	}

	return false
}

func (context *NativeRegionRenderContext) Split() {
	heap := context.Heap
	r := context.Region

	topLeftPos := r.topLeft.c
	bottomRightPos := r.bottomRight.c
	midPos := r.midPoint.c

	left := real(topLeftPos)
	right := real(bottomRightPos)
	top := imag(topLeftPos)
	bottom := imag(bottomRightPos)
	midR := real(midPos)
	midI := imag(midPos)

	topSideMid := heap.NativeEscapePoint(midR, top)
	bottomSideMid := heap.NativeEscapePoint(midR, bottom)
	leftSideMid := heap.NativeEscapePoint(left, midI)
	rightSideMid := heap.NativeEscapePoint(right, midI)

	leftSectorMid := (midR + left) / 2.0
	rightSectorMid :=  (right + midR) / 2.0
	topSectorMid := (top + midI) / 2.0
	bottomSectorMid := (midI + bottom) / 2.0

	tl := NativeRegion{
		topLeft:     r.topLeft,
		topRight:    topSideMid,
		bottomLeft:  leftSideMid,
		bottomRight: r.midPoint,
		midPoint:    heap.NativeEscapePoint(leftSectorMid, topSectorMid),
	}
	tr := NativeRegion{
		topLeft:     topSideMid,
		topRight:    r.topRight,
		bottomLeft:  r.midPoint,
		bottomRight: rightSideMid,
		midPoint:    heap.NativeEscapePoint(rightSectorMid, topSectorMid),
	}
	bl := NativeRegion{
		topLeft:     leftSideMid,
		topRight:    r.midPoint,
		bottomLeft:  r.bottomLeft,
		bottomRight: bottomSideMid,
		midPoint:    heap.NativeEscapePoint(leftSectorMid, bottomSectorMid),
	}
	br := NativeRegion{
		topLeft:     r.midPoint,
		topRight:    rightSideMid,
		bottomLeft:  bottomSideMid,
		bottomRight: r.bottomRight,
		midPoint:    heap.NativeEscapePoint(rightSectorMid, bottomSectorMid),
	}

	context.Subregion = NativeSubregion{
		populated: true,
		children:  []NativeRegion{tl, tr, bl, br},
	}
}

func (context *NativeRegionRenderContext) Rect() image.Rectangle {
	l, t := context.Config.PlaneToPixel(Nativeregion.topLeft.c)
	r, b := context.Config.PlaneToPixel(Nativeregion.bottomRight.c)
	return image.Rect(int(l), int(t), int(r), int(b))
}

func (context *NativeRegionRenderContext) Collapse() bool {
	rect := context.rect
	iCollapse := int(context.Config.NativeRegionCollapse)
	return rect.Dx() <= iCollapse || rect.Dy() <= iCollapse
}
