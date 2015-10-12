package libgodelbrot

import (
	"image"
)

type nativeSubregion struct {
	populated bool
	children  []NativeRegion
}

// Extend NativeBaseNumerics and add support for regions
type NativeRegionNumerics struct {
	Collapser
	NativeBaseNumerics
	region nativeRegion
	subregion nativeSubregion
	heap *NativeMandelbrotThunkHeap
}

type nativeRegion struct {
	topLeft     *NativeEscapePoint
	topRight    *NativeEscapePoint
	bottomLeft  *NativeEscapePoint
	bottomRight *NativeEscapePoint
	midPoint    *NativeEscapePoint
}

func (native *NativeRegionNumerics) MandelbrotPoints() {
	r := native.Region
	return []MandelbrotMember {
		r.topLeft.membership,
		r.topRight.membership,
		r.bottomLeft.membership,
		r.bottomRight.membership,
		r.midPoint.membership,
	}
}

func (native *NativeRegionNumerics) EvaluateAllPoints() {
	r := native.Region
    points := []*NativeEscapePoint{
		r.topLeft,
		r.topRight,
		r.bottomLeft,
		r.bottomRight,
		r.midPoint,
	}
    // Ensure points are all evaluated
    for _, p := range points {
        EvalThunk(p)
    }
}

// A glitch is possible when points are uniform near the set
// Due to the shape of the set, a rectangular Nativeregion is not a good approximation
// An anologous glitch happens when the entire Nativeregion is much larger than the set
// We handle both these cases here
func (native *NativeRegionNumerics) OnGlitchCurve() bool {
	member := native.RegionMember()
	iDiv := member.InvDivergence()
	iLimit, dLimit := native.MandelbrotLimits()
	if iDiv == 0 || iDiv == 1 || member.InSet() {
		sqrtChecks := native.GlitchSamples()
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
				&checkMember.Mandelbrot(iLimit, dLimit)
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

func (native *NativeRegionNumerics) Split() {
	heap := native.Heap
	r := native.Region

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

	native.Subregion = NativeSubregion{
		populated: true,
		children:  []NativeRegion{tl, tr, bl, br},
	}
}

func (native *NativeRegionNumerics) Rect() image.Rectangle {
	l, t := native.PlaneToPixel(native.Region.topLeft.c)
	r, b := native.PlaneToPixel(native.Region.bottomRight.c)
	return image.Rect(int(l), int(t), int(r), int(b))
}

// Return MandelbrotMember
// Does not check if the region's thunks have been evaluated
func (native *NativeRegionNumerics) RegionMember() MandelbrotMember {
	return native.Region.topLeft.member
}
