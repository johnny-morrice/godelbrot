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
	children  []NativeRegion
}

type nativeMandelbrotThunk struct {
	nativebase.NativeMandelbrotMember
	evaluated bool
}

type NativeRegion struct {
	topLeft     nativeMandelbrotThunk
	topRight    nativeMandelbrotThunk
	bottomLeft  nativeMandelbrotThunk
	bottomRight nativeMandelbrotThunk
	midPoint    nativeMandelbrotThunk
}

func (region *NativeRegion) rect(base *nativebase.NativeBaseNumerics) image.Rectangle {
	l, t := base.PlaneToPixel(region.topLeft.C)
	r, b := base.PlaneToPixel(region.bottomRight.C)
	return image.Rect(int(l), int(t), int(r), int(b))
}

// Extend NativeBaseNumerics and add support for regions
type NativeRegionNumerics struct {
	region.BaseRegionNumerics
	nativebase.NativeBaseNumerics
	Region             NativeRegion
	SequenceNumerics   *nativesequence.NativeSequenceNumerics
	subregion          nativeSubregion
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
			nextContexts[i] = native.Proxy(child)
		}
		return nextContexts
	}
	panic("Region asked to provide non-existent children")
	return nil
}

// Return the children of this region without hiding their types
// This implementation does not create many new objects
func (native *NativeRegionNumerics) NativeChildRegions() []NativeRegion {
	if native.subregion.populated {
		return native.subregion.children
	}
	panic("Region asked to provide non-existent children")
	return nil
}

func (native *NativeRegionNumerics) RegionSequence() region.ProxySequence {
	return native.NativeSequence()
}

func (native *NativeRegionNumerics) NativeSequence() NativeSequenceProxy {
	return NativeSequenceProxy{
		Region:   native.Region,
		NativeSequenceNumerics: native.SequenceNumerics,
	}
}

func (native *NativeRegionNumerics) Proxy(region NativeRegion) NativeRegionProxy {
	return NativeRegionProxy{
		Region:   region,
		NativeRegionNumerics: native,
	}
}

func (native *NativeRegionNumerics) MandelbrotPoints() []base.MandelbrotMember {
	r := native.Region
	return []base.MandelbrotMember{
		r.topLeft,
		r.topRight,
		r.bottomLeft,
		r.bottomRight,
		r.midPoint,
	}
}

func (native *NativeRegionNumerics) EvaluateAllPoints(iterateLimit uint8) {
	r := native.Region
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
	r := native.Region
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
	r := native.Region

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

	tl := NativeRegion{
		topLeft:     r.topLeft,
		topRight:    topSideMid,
		bottomLeft:  leftSideMid,
		bottomRight: r.midPoint,
		midPoint:    topLeftMid,
	}
	tr := NativeRegion{
		topLeft:     topSideMid,
		topRight:    r.topRight,
		bottomLeft:  r.midPoint,
		bottomRight: rightSideMid,
		midPoint:    topRightMid,
	}
	bl := NativeRegion{
		topLeft:     leftSideMid,
		topRight:    r.midPoint,
		bottomLeft:  r.bottomLeft,
		bottomRight: bottomSideMid,
		midPoint:    bottomLeftMid,
	}
	br := NativeRegion{
		topLeft:     r.midPoint,
		topRight:    rightSideMid,
		bottomLeft:  bottomSideMid,
		bottomRight: r.bottomRight,
		midPoint:    bottomRightMid,
	}

	native.subregion = nativeSubregion{
		populated: true,
		children:  []NativeRegion{tl, tr, bl, br},
	}
}

func (native *NativeRegionNumerics) Rect() image.Rectangle {
	base := native.NativeBaseNumerics
	return native.Region.rect(&base)
}

// Return MandelbrotMember
// Does not check if the region's thunks have been evaluated
func (native *NativeRegionNumerics) RegionMember() base.MandelbrotMember {
	return native.Region.topLeft
}

func createThunk(c complex128) nativeMandelbrotThunk {
	return nativeMandelbrotThunk{
		NativeMandelbrotMember: nativebase.NativeMandelbrotMember{
			C: c,
		},
	}
}