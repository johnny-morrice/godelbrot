package libgodelbrot

import (
	"image"
	"math/big"
)

type bigSubregion struct {
	populated bool
	children  []BigRegion
}

type bigMandelbrotThunk struct {
	bigMandelbrotMember
	evaluated bool
}

type bigRegion struct {
	topLeft     bigMandelbrotThunk
	topRight    bigMandelbrotThunk
	bottomLeft  bigMandelbrotThunk
	bottomRight bigMandelbrotThunk
	midPoint    bigMandelbrotThunk
}

func createBigRegion(min BigComplex, max BigComplex) bigRegion {
	left := min.Real()
	right := max.Real()
	bottom := min.Imag()
	top := max.Imag()

	midR := big.Float{}
	midR.Sub(&right, &left)
	midI := big.Float{}
	midI.Sub(&top, &bottom)

	topLeft := BigComplex{left, top}
	topRight := BigComplex{right, top}
	bottomLeft := BigComplex{left, bottom}
	bottomRight := BigComplex{right, bottom}
	midPoint := BigComplex{midR, midI}

	return bigRegion{
		topLeft:     topLeft,
		topRight:    topRight,
		bottomLeft:  bottomLeft,
		bottomRight: bottomRight,
		midPoint:    midPoint,
	}
}

// BigRegionNumerics is implementation of RegionNumerics that uses big.Float bignums for arbitrary
// accuracy.
type BigRegionNumerics struct {
	Collapser
	BigBaseNumerics
	region             bigRegion
	subregion          bigSubregion
	sequentialNumerics *BigSequentialNumerics
}

// Children returns a list of subdivided children.
func (bigFloat *BigRegionNumerics) Children() []RegionNumerics {
	if bigFloat.subregion.populated {
		nextContexts := make([]RegionNumerics, 0, 4)
		for i, child := range bigFloat.subregion.children {
			// Use a proxy to avoid heap allocation
			nextContexts[i] = bigFloat.proxyNumerics(child)
		}
		return nextContexts
	}
	panic("Region asked to provide non-existent children")
	return nil
}

// RegionalSequenceNumerics returns SequentialNumerics representing the same region on the plane.
func (bigFloat *BigRegionNumerics) RegionalSequenceNumerics() SequentialNumerics {
	return BigSequenceNumericsProxy{
		Region:   bigFloat.region,
		Numerics: bigFloat.sequentialNumerics,
	}
}

// MandelbrotPoints returns the corners of this region and its midpoint
func (bigFloat *BigRegionNumerics) MandelbrotPoints() []MandelbrotMember {
	r := bigFloat.Region
	return []MandelbrotMember{
		r.topLeft.membership,
		r.topRight.membership,
		r.bottomLeft.membership,
		r.bottomRight.membership,
		r.midPoint.membership,
	}
}

// EvaluateAllPoints runs the Mandelbrot function on all this region's points
func (bigFloat *BigRegionNumerics) EvaluateAllPoints(iterateLimit int) {
	points := []bigMandelbrotThunk{
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

// OnGlitchCurve returns true when a glitch has been detected in the region rendering process.
// The region render optimization can result in serious inaccuracies in approximating the set.
// These glitches manifest as a nasty artifact in the output image.  The cause of this is large
// areas of the plane with  same with the same InverseDivergence value.  These include the
// Mandelbrot set interior. Since regions are rectangular, we cannot well pick up on curves round
// these areas. In these cases we must perform extra sampling to be sure that we have not hit a
// glitch.
func (bigFloat *BigRegionNumerics) OnGlitchCurve(iterateLimit uint8, glitchSamples uint) bool {
	r := bigFloat.Region
	member := bigFloat.RegionMember()
	iDiv := member.InvDivergence()
	dLimit := bigFloat.DivergeLimit()
	if iDiv == 0 || iDiv == 1 || member.InSet() {
		tl := r.topLeft.c
		br := r.bottomRight.c

		hUnit := br.Real().Copy()
		hUnit.Sub(hUnit, tl.Real())
		hUnit.Quo(hUnit, glitchSamples)
		vUnit := tl.Imag().Copy()
		vUnit.Sub(h, br.Imag())
		vUnit.Quo(vUnit, glitchSamples)

		x := tl.Real()
		for i := 0; i < glitchSamples; i++ {
			y := tl.Imag()
			for j := 0; j < glitchSamples; j++ {
				checkMember := bigMandelbrotMember{
					C:            BigComplex{x, y},
					DivergeLimit: dLimit,
				}
				&checkMember.Mandelbrot(iterateLimit)
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

// Split divides the region into four smaller subregions.
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

	topSideMid := CreatebigMandelbrotThunk(midR, top)
	bottomSideMid := CreatebigMandelbrotThunk(midR, bottom)
	leftSideMid := CreatebigMandelbrotThunk(left, midI)
	rightSideMid := CreatebigMandelbrotThunk(right, midI)

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
		midPoint:    CreatebigMandelbrotThunk(leftSectorMid, topSectorMid),
	}
	tr := BigRegion{
		topLeft:     topSideMid,
		topRight:    r.topRight,
		bottomLeft:  r.midPoint,
		bottomRight: rightSideMid,
		midPoint:    CreatebigMandelbrotThunk(rightSectorMid, topSectorMid),
	}
	bl := BigRegion{
		topLeft:     leftSideMid,
		topRight:    r.midPoint,
		bottomLeft:  r.bottomLeft,
		bottomRight: bottomSideMid,
		midPoint:    CreatebigMandelbrotThunk(leftSectorMid, bottomSectorMid),
	}
	br := BigRegion{
		topLeft:     r.midPoint,
		topRight:    rightSideMid,
		bottomLeft:  bottomSideMid,
		bottomRight: r.bottomRight,
		midPoint:    CreatebigMandelbrotThunk(rightSectorMid, bottomSectorMid),
	}

	bigFloat.Subregion = BigSubregion{
		populated: true,
		children:  []BigRegion{tl, tr, bl, br},
	}
}

// Rect return a rectangle representing the position and dimensions of the region on the output
// image.
func (bigFloat *BigRegionNumerics) Rect() image.Rectangle {
	region := bigFloat.region
	l, t := bigFloat.PlaneToPixel(region.topLeft.c)
	r, b := bigFloat.PlaneToPixel(region.bottomRight.c)
	return image.Rect(int(l), int(t), int(r), int(b))
}

// RegionMember returns one MandelbrotMember from the Region
func (bigFloat *BigRegionNumerics) RegionMember() MandelbrotMember {
	return bigFloat.Region.topLeft.member.MandelbrotMember
}

// proxyNumerics quickly creates a new *NativeRegionNumerics context
func (native *NativeRegionNumerics) proxyNumerics(region Region) RegionNumerics {
	return BigRegionNumericsProxy{
		Region:   region,
		Numerics: native,
	}
}
