package bigregion

import (
	"log"
	"image"
	"math/big"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
	"functorama.com/demo/bigbase"
	"functorama.com/demo/bigsequence"
)

type bigSubregion struct {
	populated bool
	children  []bigRegion
}

type bigMandelbrotThunk struct {
	bigbase.BigMandelbrotMember
	evaluated bool
	cStore bigbase.BigComplex
}

type bigRegion struct {
	topLeft     bigMandelbrotThunk
	topRight    bigMandelbrotThunk
	bottomLeft  bigMandelbrotThunk
	bottomRight bigMandelbrotThunk
	midPoint    bigMandelbrotThunk
}

func (br *bigRegion) thunks() []*bigMandelbrotThunk {
	return []*bigMandelbrotThunk{
		&br.topLeft,
		&br.topRight,
		&br.bottomLeft,
		&br.bottomRight,
		&br.midPoint,
	}
}

// Rect return a rectangle representing the position and dimensions of the region on the output
// image.
func (br *bigRegion) rect(base *bigbase.BigBaseNumerics) image.Rectangle {
	l, t := base.PlaneToPixel(br.topLeft.C)
	r, b := base.PlaneToPixel(br.bottomRight.C)
	return image.Rect(l, t, r, b)
}

func (brn *BigRegionNumerics) createRelativeRegion(min *bigbase.BigComplex, max *bigbase.BigComplex) bigRegion {
	left := min.Real()
	right := max.Real()
	bottom := min.Imag()
	top := max.Imag()

	midR := brn.MakeBigFloat(0.0)
	midR.Sub(right, left)
	midI := brn.MakeBigFloat(0.0)
	midI.Sub(top, bottom)

	return bigRegion{
		topLeft:     brn.createBigThunk(left, top),
		topRight:    brn.createBigThunk(right, top),
		bottomLeft:  brn.createBigThunk(bottom, left),
		bottomRight: brn.createBigThunk(bottom, right),
		midPoint:    brn.createBigThunk(&midR, &midI),
	}
}

func (brn *BigRegionNumerics) createBigThunk(r *big.Float, i *big.Float) bigMandelbrotThunk {
	thunk := bigMandelbrotThunk{
		cStore: bigbase.BigComplex{*r, *i},
	}
	thunk.BigMandelbrotMember.C = &thunk.cStore
	thunk.BigMandelbrotMember.SqrtDivergeLimit = &brn.SqrtDivergeLimit
	return thunk
}

// BigRegionNumerics is implementation of RegionNumerics that uses big.Float bignums for arbitrary
// accuracy.
type BigRegionNumerics struct {
	region.RegionConfig
	bigbase.BigBaseNumerics
	Region             bigRegion
	SequenceNumerics  *bigsequence.BigSequenceNumerics
	subregion          bigSubregion
}
var _ region.RegionNumerics = (*BigRegionNumerics)(nil)

func CreateBigRegionNumerics(app RenderApplication) BigRegionNumerics {
	sequence := bigsequence.CreateBigSequenceNumerics(app)
	parent := bigbase.CreateBigBaseNumerics(app)
	planeMin := bigbase.BigComplex{parent.RealMin, parent.ImagMin}
	planeMax := bigbase.BigComplex{parent.RealMax, parent.ImagMax}
	return BigRegionNumerics{
		BigBaseNumerics: parent,
		RegionConfig: app.RegionConfig(),
		SequenceNumerics: &sequence,
		Region: createBigRegion(planeMin, planeMax),
	}
}

func (brn *BigRegionNumerics) ClaimExtrinsics() {
	// We have our extrinsics right here
}

func (brn *BigRegionNumerics) BigChildRegions() []bigRegion {
	if brn.subregion.populated {
		return brn.subregion.children
	}
	log.Panic("No children!")
	return nil
}

// Children returns a list of subdivided children.
func (brn *BigRegionNumerics) Children() []region.RegionNumerics {
	if brn.subregion.populated {
		nextContexts := make([]region.RegionNumerics, 4)
		for i, child := range brn.subregion.children {
			// Use a proxy to avoid heap allocation
			nextContexts[i] = brn.proxyNumerics(&child)
		}
		return nextContexts
	}
	log.Panic("Region asked to provide non-existent children")
	return nil
}

// RegionSequence returns ProxySequenceNumerics representing the same region on the plane.
func (brn *BigRegionNumerics) RegionSequence() region.ProxySequence {
	return BigSequenceNumericsProxy{
		BigSequenceNumerics: brn.SequenceNumerics,
		LocalRegion:   brn.Region,
	}
}

// MandelbrotPoints returns the corners of this region and its midpoint
func (brn *BigRegionNumerics) MandelbrotPoints() []base.MandelbrotMember {
	r := brn.Region
	return []base.MandelbrotMember{
		r.topLeft,
		r.topRight,
		r.bottomLeft,
		r.bottomRight,
		r.midPoint,
	}
}

// EvaluateAllPoints runs the Mandelbrot function on all this region's points
func (brn *BigRegionNumerics) EvaluateAllPoints(iterateLimit uint8) {
	thunks := brn.Region.thunks()

	// Ensure points are all evaluated
	for _, p := range thunks {
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
func (brn *BigRegionNumerics) OnGlitchCurve(iterateLimit uint8, glitchSamples uint) bool {
	r := brn.Region
	member := brn.Region.topLeft
	iDiv := member.InvDivergence
	if iDiv == 0 || iDiv == 1 || member.InSet {
		tl := r.topLeft.cStore
		br := r.bottomRight.cStore

		bigGlitch := brn.MakeBigFloat(float64(glitchSamples))

		rUnit := brn.MakeBigFloat(0.0)
		rUnit.Copy(br.Real())
		rUnit.Sub(&rUnit, tl.Real())
		rUnit.Quo(&rUnit, &bigGlitch)
		iUnit := brn.MakeBigFloat(0.0)
		iUnit.Copy(tl.Imag())
		iUnit.Sub(&iUnit, br.Imag())
		iUnit.Quo(&iUnit, &bigGlitch)

		startI := *tl.Imag()
		c := bigbase.BigComplex{R: *tl.Real()}
		for i := uint(0); i < glitchSamples; i++ {
			c.I = startI
			for j := uint(0); j < glitchSamples; j++ {
				checkMember := bigbase.BigMandelbrotMember{
					C:            &c,
					SqrtDivergeLimit: &brn.SqrtDivergeLimit,
				}
				checkMember.Mandelbrot(iterateLimit)
				if member.InvDivergence != iDiv {
					return true
				}
				c.I.Sub(&c.I, &iUnit)
			}
			c.R.Add(&c.R, &rUnit)
		}
	}

	return false
}

// Split divides the region into four smaller subregions.
func (brn *BigRegionNumerics) Split() {
	r := brn.Region

	topLeftPos := r.topLeft.cStore
	bottomRightPos := r.bottomRight.cStore
	midPos := r.midPoint.cStore

	left := topLeftPos.Real()
	right := bottomRightPos.Real()
	top := topLeftPos.Imag()
	bottom := bottomRightPos.Imag()
	midR := midPos.Real()
	midI := midPos.Imag()

	bigTwo := brn.MakeBigFloat(2.0)

	topSideMid := brn.createBigThunk(midR, top)
	bottomSideMid := brn.createBigThunk(midR, bottom)
	leftSideMid := brn.createBigThunk(left, midI)
	rightSideMid := brn.createBigThunk(right, midI)

	leftSectorMid := brn.MakeBigFloat(0.0)
	leftSectorMid.Copy(left)
	leftSectorMid.Add(&leftSectorMid, midR)
	leftSectorMid.Quo(&leftSectorMid, &bigTwo)

	rightSectorMid := brn.MakeBigFloat(0.0)
	rightSectorMid.Copy(right)
	rightSectorMid.Add(&rightSectorMid, midR)
	rightSectorMid.Quo(&rightSectorMid, &bigTwo)

	topSectorMid := brn.MakeBigFloat(0.0)
	topSectorMid.Copy(top)
	topSectorMid.Add(&topSectorMid, midI)
	topSectorMid.Quo(&topSectorMid, &bigTwo)

	bottomSectorMid := brn.MakeBigFloat(0.0)
	bottomSectorMid.Copy(bottom)
	bottomSectorMid.Add(&bottomSectorMid, midI)
	bottomSectorMid.Quo(&bottomSectorMid, &bigTwo)

	tl := bigRegion{
		topLeft:     r.topLeft,
		topRight:    topSideMid,
		bottomLeft:  leftSideMid,
		bottomRight: r.midPoint,
		midPoint:    brn.createBigThunk(&leftSectorMid, &topSectorMid),
	}
	tr := bigRegion{
		topLeft:     topSideMid,
		topRight:    r.topRight,
		bottomLeft:  r.midPoint,
		bottomRight: rightSideMid,
		midPoint:    brn.createBigThunk(&rightSectorMid, &topSectorMid),
	}
	bl := bigRegion{
		topLeft:     leftSideMid,
		topRight:    r.midPoint,
		bottomLeft:  r.bottomLeft,
		bottomRight: bottomSideMid,
		midPoint:    brn.createBigThunk(&leftSectorMid, &bottomSectorMid),
	}
	br := bigRegion{
		topLeft:     r.midPoint,
		topRight:    rightSideMid,
		bottomLeft:  bottomSideMid,
		bottomRight: r.bottomRight,
		midPoint:    brn.createBigThunk(&rightSectorMid, &bottomSectorMid),
	}

	brn.subregion = bigSubregion{
		populated: true,
		children:  []bigRegion{tl, tr, bl, br},
	}
}

// RegionMember returns one MandelbrotMember from the Region
func (brn *BigRegionNumerics) RegionMember() base.MandelbrotMember {
	return brn.Region.topLeft.BigMandelbrotMember
}

// proxyNumerics quickly creates a new *NativeRegionNumerics context
func (brn *BigRegionNumerics) proxyNumerics(region *bigRegion) region.RegionNumerics {
	return BigRegionNumericsProxy{
		BigRegionNumerics: brn,
		LocalRegion:   *region,
	}
}

func (brn *BigRegionNumerics) Rect() image.Rectangle {
	return brn.Region.rect(&brn.BigBaseNumerics)
}

func createBigRegion(min bigbase.BigComplex, max bigbase.BigComplex) bigRegion {
	left := min.R
	bottom := min.I
	right := max.R
	top := max.I

	prec := left.Prec()
	bigTwo := bigbase.MakeBigFloat(2.0, prec)

	midR := bigbase.MakeBigFloat(0.0, prec)
	midR.Add(&right, &left)
	midR.Quo(&midR, &bigTwo)

	midI := bigbase.MakeBigFloat(0.0, prec)
	midI.Add(&top, &bottom)
	midI.Quo(&midI, &bigTwo)

	corners := []bigbase.BigComplex{
		bigbase.BigComplex{left, top},
		bigbase.BigComplex{right, top},
		bigbase.BigComplex{left, bottom},
		bigbase.BigComplex{right, bottom},
		bigbase.BigComplex{midR, midI}, // midpoint is not technically a corner
	}

	thunks := make([]bigMandelbrotThunk, len(corners))

	for i, c := range corners {
		thunks[i] = bigMandelbrotThunk{
			cStore: c,
		}
		thunks[i].C = &thunks[i].cStore
	}

	return bigRegion{
		topLeft: thunks[0],
		topRight: thunks[1],
		bottomLeft: thunks[2],
		bottomRight: thunks[3],
		midPoint: thunks[4],
	}
}