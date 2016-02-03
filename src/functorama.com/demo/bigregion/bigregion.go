package bigregion

import (
	"log"
	"image"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
	"functorama.com/demo/bigbase"
	"functorama.com/demo/bigsequence"
)

type bigSubregion struct {
	populated bool
	children  []bigRegion
}

type bigRegion struct {
	topLeft     bigbase.BigMandelbrotMember
	topRight    bigbase.BigMandelbrotMember
	bottomLeft  bigbase.BigMandelbrotMember
	bottomRight bigbase.BigMandelbrotMember
	midPoint    bigbase.BigMandelbrotMember
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
		topLeft:     brn.MakeMember(&bigbase.BigComplex{*left, *top}),
		topRight:    brn.MakeMember(&bigbase.BigComplex{*right, *top}),
		bottomLeft:  brn.MakeMember(&bigbase.BigComplex{*bottom, *left}),
		bottomRight: brn.MakeMember(&bigbase.BigComplex{*bottom, *right}),
		midPoint:    brn.MakeMember(&bigbase.BigComplex{midR, midI}),
	}
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

func Make(app RenderApplication) BigRegionNumerics {
	sequence := bigsequence.Make(app)
	parent := bigbase.Make(app)
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

func (brn *BigRegionNumerics) Extrinsically(f func()) {
	f()
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
	ps := brn.Points()
	base := make([]base.MandelbrotMember, len(ps))
	for i, p := range ps {
		base[i] = p.MandelbrotMember
	}
	return base
}

func (brn *BigRegionNumerics) Points() []bigbase.BigMandelbrotMember {
	r := brn.Region
	return []bigbase.BigMandelbrotMember{
		r.topLeft,
		r.topRight,
		r.midPoint,
		r.bottomLeft,
		r.bottomRight,
	}
}

// Split divides the region into four smaller subregions.
func (brn *BigRegionNumerics) Split() {
	r := brn.Region

	topLeftPos := r.topLeft.C
	bottomRightPos := r.bottomRight.C
	midPos := r.midPoint.C

	left := topLeftPos.Real()
	right := bottomRightPos.Real()
	top := topLeftPos.Imag()
	bottom := bottomRightPos.Imag()
	midR := midPos.Real()
	midI := midPos.Imag()

	bigTwo := brn.MakeBigFloat(2.0)

	topSideMid := brn.MakeMember(&bigbase.BigComplex{*midR, *top})
	bottomSideMid := brn.MakeMember(&bigbase.BigComplex{*midR, *bottom})
	leftSideMid := brn.MakeMember(&bigbase.BigComplex{*left, *midI})
	rightSideMid := brn.MakeMember(&bigbase.BigComplex{*right, *midI})

	leftSectorMid := brn.MakeBigFloat(0.0)
	leftSectorMid.Add(left, midR)
	leftSectorMid.Quo(&leftSectorMid, &bigTwo)

	rightSectorMid := brn.MakeBigFloat(0.0)
	rightSectorMid.Add(right, midR)
	rightSectorMid.Quo(&rightSectorMid, &bigTwo)

	topSectorMid := brn.MakeBigFloat(0.0)
	topSectorMid.Add(top, midI)
	topSectorMid.Quo(&topSectorMid, &bigTwo)

	bottomSectorMid := brn.MakeBigFloat(0.0)
	bottomSectorMid.Add(bottom, midI)
	bottomSectorMid.Quo(&bottomSectorMid, &bigTwo)

	tl := bigRegion{
		topLeft:     r.topLeft,
		topRight:    topSideMid,
		bottomLeft:  leftSideMid,
		bottomRight: r.midPoint,
		midPoint:    brn.MakeMember(&bigbase.BigComplex{leftSectorMid, topSectorMid}),
	}
	tr := bigRegion{
		topLeft:     topSideMid,
		topRight:    r.topRight,
		bottomLeft:  r.midPoint,
		bottomRight: rightSideMid,
		midPoint:    brn.MakeMember(&bigbase.BigComplex{rightSectorMid, topSectorMid}),
	}
	bl := bigRegion{
		topLeft:     leftSideMid,
		topRight:    r.midPoint,
		bottomLeft:  r.bottomLeft,
		bottomRight: bottomSideMid,
		midPoint:    brn.MakeMember(&bigbase.BigComplex{leftSectorMid, bottomSectorMid}),
	}
	br := bigRegion{
		topLeft:     r.midPoint,
		topRight:    rightSideMid,
		bottomLeft:  bottomSideMid,
		bottomRight: r.bottomRight,
		midPoint:    brn.MakeMember(&bigbase.BigComplex{rightSectorMid, bottomSectorMid}),
	}

	brn.subregion = bigSubregion{
		populated: true,
		children:  []bigRegion{tl, tr, bl, br},
	}
}

// RegionMember returns one MandelbrotMember from the Region
func (brn *BigRegionNumerics) RegionMember() base.MandelbrotMember {
	return brn.Region.topLeft.MandelbrotMember
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

// TODO
func (brn *BigRegionNumerics) SampleDivs() (<-chan uint8, chan<- bool) {
	done := make(chan bool, 1)
	idivch := make(chan uint8, 1)

	return idivch, done
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

	thunks := make([]bigbase.BigMandelbrotMember, len(corners))

	for i, c := range corners {
		thunks[i] = bigbase.BigMandelbrotMember{
			C: &c,
		}
	}

	return bigRegion{
		topLeft: thunks[0],
		topRight: thunks[1],
		bottomLeft: thunks[2],
		bottomRight: thunks[3],
		midPoint: thunks[4],
	}
}