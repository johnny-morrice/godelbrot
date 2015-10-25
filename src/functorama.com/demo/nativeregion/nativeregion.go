package nativeregion

import (
	"image"
	"functorama.com/demo/base"
	"functorama.com/demo/nativebase"
	"functorama.com/demo/nativesequence"
	"functorama.com/demo/region"
)

type nativeSubregion struct {
	populated bool
	children  []nativeRegion
}

type nativeMandelbrotThunk struct {
	nativebase.NativeMandelbrotMember
	evaluated bool
}

type nativeRegion struct {
	topLeft     nativeMandelbrotThunk
	topRight    nativeMandelbrotThunk
	bottomLeft  nativeMandelbrotThunk
	bottomRight nativeMandelbrotThunk
	midPoint    nativeMandelbrotThunk
}

func (region *nativeRegion) rect(base *nativebase.NativeBaseNumerics) image.Rectangle {
	l, t := base.PlaneToPixel(region.topLeft.C)
	r, b := base.PlaneToPixel(region.bottomRight.C)
	return image.Rect(int(l), int(t), int(r), int(b))
}

// Extend NativeBaseNumerics and add support for regions
type NativeRegionNumerics struct {
	base.BaseRegionNumerics
	nativebase.NativeBaseNumerics
	region             nativeRegion
	subregion          nativeSubregion
	sequenceNumerics   *nativesequence.NativeSequenceNumerics
}

func (native *NativeRegionNumerics) ClaimExtrinsics() {
	// Region already present
}

// Return the children of this region
// This implementation does not create many new objects
func (native *NativeRegionNumerics) Children() []region.RegionNumerics {
	if native.subregion.populated {
		nextContexts := make([]region.RegionNumerics, 0, 4)
		for i, child := range native.subregion.children {
			nextContexts[i] = native.proxyNumerics(child)
		}
		return nextContexts
	}
	panic("Region asked to provide non-existent children")
	return nil
}

func (native *NativeRegionNumerics) RegionSequenceNumerics() region.RegionSequenceNumerics {
	return NativeSequenceNumericsProxy{
		region:   native.region,
		NativeSequenceNumerics: native.sequenceNumerics,
	}
}

func (native *NativeRegionNumerics) MandelbrotPoints() []base.MandelbrotMember {
	r := native.region
	return []base.MandelbrotMember{
		r.topLeft,
		r.topRight,
		r.bottomLeft,
		r.bottomRight,
		r.midPoint,
	}
}

func (native *NativeRegionNumerics) EvaluateAllPoints(iterateLimit uint8) {
	r := native.region
	points := []nativeMandelbrotThunk{
		r.topLeft,
		r.topRight,
		r.bottomLeft,
		r.bottomRight,
		r.midPoint,
	}
	// Ensure points are all evaluated
	for _, p := range points {
		if !p.evaluated {
			p.Mandelbrot(iterateLimit)
			p.evaluated = true
		}
	}
}

// A glitch is possible when points are uniform near the set
// Due to the shape of the set, a rectangular Nativeregion is not a good approximation
// An anologous glitch happens when the entire Nativeregion is much larger than the set
// We handle both these cases here
func (native *NativeRegionNumerics) OnGlitchCurve(iterateLimit uint8, glitchSamples uint) bool {
	r := native.region
	tlMember := r.topLeft
	iDiv := tlMember.InvDivergence
	if iDiv == 0 || iDiv == 1 || tlMember.InSet {
		sqrtDLimit := native.SqrtDivergeLimit
		glitchSamplesF := float64(glitchSamples)
		tl := tlMember.C
		br := r.bottomRight.C
		w := real(br) - real(tl)
		h := imag(tl) - imag(br)
		vUnit := h / glitchSamplesF
		hUnit := w / glitchSamplesF
		x := real(tl)
		for i := uint(0); i < glitchSamples; i++ {
			y := imag(tl)
			for j := uint(0); j < glitchSamples; j++ {
				checkMember := nativebase.NativeMandelbrotMember{
					C: complex(x, y),
					SqrtDivergeLimit: sqrtDLimit,
				}
				checkMember.Mandelbrot(iterateLimit)
				if checkMember.InvDivergence != iDiv {
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
	r := native.region

	topLeftPos := r.topLeft.C
	bottomRightPos := r.bottomRight.C
	midPos := r.midPoint.C

	left := real(topLeftPos)
	right := real(bottomRightPos)
	top := imag(topLeftPos)
	bottom := imag(bottomRightPos)
	midR := real(midPos)
	midI := imag(midPos)

	leftSectorMid := (midR + left) / 2.0
	rightSectorMid := (right + midR) / 2.0
	topSectorMid := (top + midI) / 2.0
	bottomSectorMid := (midI + bottom) / 2.0

	topSideMid := createThunk(complex(midR, top))
	bottomSideMid := createThunk(complex(midR, bottom))
	leftSideMid := createThunk(complex(left, midI))
	rightSideMid := createThunk(complex(right, midI))

	topLeftMid := createThunk(complex(leftSectorMid, topSectorMid))
	topRightMid := createThunk(complex(rightSectorMid, topSectorMid))
	bottomLeftMid := createThunk(complex(leftSectorMid, bottomSectorMid))
	bottomRightMid := createThunk(complex(rightSectorMid, bottomSectorMid))

	tl := nativeRegion{
		topLeft:     r.topLeft,
		topRight:    topSideMid,
		bottomLeft:  leftSideMid,
		bottomRight: r.midPoint,
		midPoint:    topLeftMid,
	}
	tr := nativeRegion{
		topLeft:     topSideMid,
		topRight:    r.topRight,
		bottomLeft:  r.midPoint,
		bottomRight: rightSideMid,
		midPoint:    topRightMid,
	}
	bl := nativeRegion{
		topLeft:     leftSideMid,
		topRight:    r.midPoint,
		bottomLeft:  r.bottomLeft,
		bottomRight: bottomSideMid,
		midPoint:    bottomLeftMid,
	}
	br := nativeRegion{
		topLeft:     r.midPoint,
		topRight:    rightSideMid,
		bottomLeft:  bottomSideMid,
		bottomRight: r.bottomRight,
		midPoint:    bottomRightMid,
	}

	native.subregion = nativeSubregion{
		populated: true,
		children:  []nativeRegion{tl, tr, bl, br},
	}
}

func (native *NativeRegionNumerics) Rect() image.Rectangle {
	base := native.NativeBaseNumerics
	return native.region.rect(&base)
}

// Return MandelbrotMember
// Does not check if the region's thunks have been evaluated
func (native *NativeRegionNumerics) RegionMember() base.MandelbrotMember {
	return native.region.topLeft
}

// Quickly create a new *NativeRegionNumerics context
func (native *NativeRegionNumerics) proxyNumerics(region nativeRegion) region.RegionNumerics {
	return NativeRegionNumericsProxy{
		region:   region,
		NativeRegionNumerics: native,
	}
}

func createThunk(c complex128) nativeMandelbrotThunk {
	return nativeMandelbrotThunk{
		NativeMandelbrotMember: nativebase.NativeMandelbrotMember{
			C: c,	
		},
	}
}