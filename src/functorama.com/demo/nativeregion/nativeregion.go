package nativeregion

import (
	"log"
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

func (nr *nativeRegion) rect(base *nativebase.NativeBaseNumerics) image.Rectangle {
	l, t := base.PlaneToPixel(nr.topLeft.C)
	r, b := base.PlaneToPixel(nr.bottomRight.C)
	return image.Rect(l, t, r, b)
}

func (nr *nativeRegion) thunks() []*nativeMandelbrotThunk {
	return []*nativeMandelbrotThunk{
		&nr.topLeft,
		&nr.topRight,
		&nr.bottomLeft,
		&nr.bottomRight,
		&nr.midPoint,
	}
}

// Extend NativeBaseNumerics and add support for regions
type NativeRegionNumerics struct {
	region.RegionConfig
	nativebase.NativeBaseNumerics
	Region             nativeRegion
	SequenceNumerics   *nativesequence.NativeSequenceNumerics
	subregion          nativeSubregion
}

// Check that we implement the interface
var _ region.RegionNumerics = (*NativeRegionNumerics)(nil)

func Make(app RenderApplication) NativeRegionNumerics {
	sequence := nativesequence.CreateNativeSequenceNumerics(app)
	parent := nativebase.Make(app)
	planeMin := complex(parent.RealMin, parent.ImagMin)
	planeMax := complex(parent.RealMax, parent.ImagMax)
	return NativeRegionNumerics{
		NativeBaseNumerics: parent,
		RegionConfig: app.RegionConfig(),
		SequenceNumerics: &sequence,
		Region: createNativeRegion(planeMin, planeMax),
	}
}

func (native *NativeRegionNumerics) ClaimExtrinsics() {
	// Region already present
}

// Return the children of this region
// This implementation does not create many new objects
func (native *NativeRegionNumerics) Children() []region.RegionNumerics {
	const childCount = 4
	if native.subregion.populated {
		nextContexts := make([]region.RegionNumerics, childCount)
		for i, child := range native.subregion.children {
			nextContexts[i] = native.Proxy(child)
		}
		return nextContexts
	}
	log.Panic("Region asked to provide non-existent children")
	return nil
}

// Return the children of this region without hiding their types
// This implementation does not create many new objects
func (native *NativeRegionNumerics) NativeChildRegions() []nativeRegion {
	if native.subregion.populated {
		return native.subregion.children
	}
	log.Panic("Region asked to provide non-existent children")
	return nil
}

func (native *NativeRegionNumerics) RegionSequence() region.ProxySequence {
	return native.NativeSequence()
}

func (native *NativeRegionNumerics) NativeSequence() NativeSequenceProxy {
	return NativeSequenceProxy{
		LocalRegion:   native.Region,
		NativeSequenceNumerics: native.SequenceNumerics,
	}
}

func (native *NativeRegionNumerics) Proxy(region nativeRegion) NativeRegionProxy {
	return NativeRegionProxy{
		LocalRegion:   region,
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
	points := native.Region.thunks()
	// Ensure points are all evaluated
	for _, p := range points {
		if !p.evaluated {
			p.Mandelbrot(iterateLimit)
			p.evaluated = true
		}
	}
}

// A glitch is possible when points are uniform near the set
// Due to the shape of the set, a rectangular nativeRegion is not a good approximation
// An anologous glitch happens when the entire nativeRegion is much larger than the set
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

	topSideMid := native.createThunk(complex(midR, top))
	bottomSideMid := native.createThunk(complex(midR, bottom))
	leftSideMid := native.createThunk(complex(left, midI))
	rightSideMid := native.createThunk(complex(right, midI))

	topLeftMid := native.createThunk(complex(leftSectorMid, topSectorMid))
	topRightMid := native.createThunk(complex(rightSectorMid, topSectorMid))
	bottomLeftMid := native.createThunk(complex(leftSectorMid, bottomSectorMid))
	bottomRightMid := native.createThunk(complex(rightSectorMid, bottomSectorMid))

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
	return native.Region.rect(&base)
}

// Return MandelbrotMember
// Does not check if the region's thunks have been evaluated
func (native *NativeRegionNumerics) RegionMember() base.MandelbrotMember {
	return native.Region.topLeft
}

func (native *NativeRegionNumerics) createThunk(c complex128) nativeMandelbrotThunk {
	return nativeMandelbrotThunk{
		NativeMandelbrotMember: native.CreateMandelbrot(c),
	}
}

func (native *NativeRegionNumerics) thunks() []nativeMandelbrotThunk {
	region := native.Region
	return []nativeMandelbrotThunk{
		region.topLeft,
		region.topRight,
		region.bottomLeft,
		region.bottomRight,
		region.midPoint,
	}
}

func createNativeRegion(min complex128, max complex128) nativeRegion {
	left := real(min)
	right := real(max)
	top := imag(max)
	bottom := imag(min)
	mid := ((max - min) / 2) + min

	points := []complex128{
		complex(left, top),
		complex(right, top),
		complex(left, bottom),
		complex(right, bottom),
		mid,
	}

	thunks := make([]nativeMandelbrotThunk, len(points))
	for i, c := range points {
		thunks[i] = nativeMandelbrotThunk{
			NativeMandelbrotMember: nativebase.NativeMandelbrotMember{C: c, SqrtDivergeLimit: 2.0},
		}
	}

	region := nativeRegion{
		topLeft: thunks[0],
		topRight: thunks[1],
		bottomLeft: thunks[2],
		bottomRight: thunks[3],
		midPoint: thunks[4],
	}

	return region
}