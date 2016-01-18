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

type nativeRegion struct {
	topLeft     nativebase.NativeMandelbrotMember
	topRight    nativebase.NativeMandelbrotMember
	bottomLeft  nativebase.NativeMandelbrotMember
	bottomRight nativebase.NativeMandelbrotMember
	midPoint    nativebase.NativeMandelbrotMember
}

func (nr *nativeRegion) rect(base *nativebase.NativeBaseNumerics) image.Rectangle {
	l, t := base.PlaneToPixel(nr.topLeft.C)
	r, b := base.PlaneToPixel(nr.bottomRight.C)
	return image.Rect(l, t, r, b)
}

func (nr *nativeRegion) points() []*nativebase.NativeMandelbrotMember {
	return []*nativebase.NativeMandelbrotMember{
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
	sequence := nativesequence.Make(app)
	parent := nativebase.Make(app)
	planeMin := complex(parent.RealMin, parent.ImagMin)
	planeMax := complex(parent.RealMax, parent.ImagMax)
	divergeLimit := parent.SqrtDivergeLimit
	return NativeRegionNumerics{
		NativeBaseNumerics: parent,
		RegionConfig: app.RegionConfig(),
		SequenceNumerics: &sequence,
		Region: createNativeRegion(planeMin, planeMax, divergeLimit),
	}
}

func (native *NativeRegionNumerics) ClaimExtrinsics() {
	// Region already present
}

func (native *NativeRegionNumerics) Extrinsically(f func()) {
	f()
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
	seq := NativeSequenceProxy{
		LocalRegion:   native.Region,
		NativeSequenceNumerics: native.SequenceNumerics,
	}
	return seq
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

func (native *NativeRegionNumerics) Split(iterateLimit uint8) {
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

	topSideMid := native.createPoint(complex(midR, top), iterateLimit)
	bottomSideMid := native.createPoint(complex(midR, bottom), iterateLimit)
	leftSideMid := native.createPoint(complex(left, midI), iterateLimit)
	rightSideMid := native.createPoint(complex(right, midI), iterateLimit)

	topLeftMid := native.createPoint(complex(leftSectorMid, topSectorMid), iterateLimit)
	topRightMid := native.createPoint(complex(rightSectorMid, topSectorMid), iterateLimit)
	bottomLeftMid := native.createPoint(complex(leftSectorMid, bottomSectorMid), iterateLimit)
	bottomRightMid := native.createPoint(complex(rightSectorMid, bottomSectorMid), iterateLimit)

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
// Does not check if the region's Points have been evaluated
func (native *NativeRegionNumerics) RegionMember() base.MandelbrotMember {
	return native.Region.topLeft
}

func (native *NativeRegionNumerics) createPoint(c complex128, iterateLimit uint8) nativebase.NativeMandelbrotMember {
	point := native.CreateMandelbrot(c)
	point.Mandelbrot(iterateLimit)
	return point
}

func (native *NativeRegionNumerics) Points() []nativebase.NativeMandelbrotMember {
	region := native.Region
	return []nativebase.NativeMandelbrotMember{
		region.topLeft,
		region.topRight,
		region.bottomLeft,
		region.bottomRight,
		region.midPoint,
	}
}

func createNativeRegion(min complex128, max complex128, sqrtDLimit float64) nativeRegion {
	left := real(min)
	right := real(max)
	top := imag(max)
	bottom := imag(min)
	mid := ((max - min) / 2) + min

	coords := []complex128{
		complex(left, top),
		complex(right, top),
		complex(left, bottom),
		complex(right, bottom),
		mid,
	}

	points := make([]nativebase.NativeMandelbrotMember, len(coords))
	for i, c := range coords {
		points[i] = nativebase.NativeMandelbrotMember{C: c, SqrtDivergeLimit: sqrtDLimit}
	}

	region := nativeRegion{
		topLeft: points[0],
		topRight: points[1],
		bottomLeft: points[2],
		bottomRight: points[3],
		midPoint: points[4],
	}

	return region
}