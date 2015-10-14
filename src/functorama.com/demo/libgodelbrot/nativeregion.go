package libgodelbrot

import (
	"image"
)

type nativeSubregion struct {
	populated bool
	children  []NativeRegion
}

type nativeRegion struct {
	topLeft     NativeMandelbrotMember
	topRight    NativeMandelbrotMember
	bottomLeft  NativeMandelbrotMember
	bottomRight NativeMandelbrotMember
	midPoint    NativeMandelbrotMember
}

// Extend NativeBaseNumerics and add support for regions
type NativeRegionNumerics struct {
	BaseRegionNumerics
	NativeBaseNumerics
	region nativeRegion
	subregion nativeSubregion
	sequentialNumerics *NativeSequentialNumerics
}

// Return the children of this region
// This implementation does not create many new objects
func (native *NativeRegionNumerics) Children() []RegionNumerics {
	if native.subregion.populated {
		nextContexts := make([]RegionNumerics, 0, 4)
		for i, child := range native.subregion.children {
			nextContexts[i] := native.proxyNumerics(child)
		}
		return nextContexts
	}
	panic("Region asked to provide non-existent children")
	return nil
}

func (native *NativeRegionNumerics) RegionalSequenceNumerics() {
	return NativeSequenceNumericsProxy{
		Region: native.region,
		Numerics: native.sequentialNumerics,
	}
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
    points := []NativeMandelbrotMember{
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

	topSideMid := NativeMandelbrotMember{C: complex(midR, top)}
	bottomSideMid := NativeMandelbrotMember{C: complex(midR, bottom)}
	leftSideMid := NativeMandelbrotMember{C: complex(left, midI)}
	rightSideMid := NativeMandelbrotMember{C: complex(right, midI)}

	leftSectorMid := (midR + left) / 2.0
	rightSectorMid :=  (right + midR) / 2.0
	topSectorMid := (top + midI) / 2.0
	bottomSectorMid := (midI + bottom) / 2.0

	tl := NativeRegion{
		topLeft:     r.topLeft,
		topRight:    topSideMid,
		bottomLeft:  leftSideMid,
		bottomRight: r.midPoint,
		midPoint:    NativeMandelbrotMember{C: complex(leftSectorMid, topSectorMid)},
	}
	tr := NativeRegion{
		topLeft:     topSideMid,
		topRight:    r.topRight,
		bottomLeft:  r.midPoint,
		bottomRight: rightSideMid,
		midPoint:    NativeMandelbrotMember{C: complex(rightSectorMid, topSectorMid)},
	}
	bl := NativeRegion{
		topLeft:     leftSideMid,
		topRight:    r.midPoint,
		bottomLeft:  r.bottomLeft,
		bottomRight: bottomSideMid,
		midPoint:    NativeMandelbrotMember{C: complex(leftSectorMid, bottomSectorMid)},
	}
	br := NativeRegion{
		topLeft:     r.midPoint,
		topRight:    rightSideMid,
		bottomLeft:  bottomSideMid,
		bottomRight: r.bottomRight,
		midPoint:    NativeMandelbrotMember{C: complex(rightSectorMid, bottomSectorMid)},
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

// Quickly create a new *NativeRegionNumerics context
func (native *NativeRegionNumerics) proxyNumerics(region Region) RegionNumerics {
	return NativeRegionNumericsProxy{
		Region: region,
		Numerics: native,
	}
}