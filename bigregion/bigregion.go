package bigregion

import (
	"log"
	"image"
	"math/big"
	"github.com/johnny-morrice/godelbrot/base"
	"github.com/johnny-morrice/godelbrot/region"
	"github.com/johnny-morrice/godelbrot/bigbase"
	"github.com/johnny-morrice/godelbrot/bigsequence"
)

type bigSubregion struct {
	populated bool
	children  []bigRegion
}

type bigRegion struct {
	topLeft     bigbase.BigEscapeValue
	topRight    bigbase.BigEscapeValue
	bottomLeft  bigbase.BigEscapeValue
	bottomRight bigbase.BigEscapeValue
	midPoint    bigbase.BigEscapeValue
}

func (br *bigRegion) points() []bigbase.BigEscapeValue {
	return []bigbase.BigEscapeValue{
		br.topLeft,
		br.topRight,
		br.midPoint,
		br.bottomLeft,
		br.bottomRight,
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
func (brn *BigRegionNumerics) MandelbrotPoints() []base.EscapeValue {
	ps := brn.Points()
	base := make([]base.EscapeValue, len(ps))
	for i, p := range ps {
		base[i] = p.EscapeValue
	}
	return base
}

func (brn *BigRegionNumerics) Points() []bigbase.BigEscapeValue {
	return brn.Region.points()
}

// Split divides the region into four smaller subregions.
func (brn *BigRegionNumerics) Split() {
	r := brn.Region

	topLeftPos := r.topLeft.C
	bottomRightPos := r.bottomRight.C
	midPos := r.midPoint.C

	left := topLeftPos.R
	right := bottomRightPos.R
	top := topLeftPos.I
	bottom := bottomRightPos.I
	midR := midPos.R
	midI := midPos.I

	bigTwo := brn.MakeBigFloat(2.0)

	leftSectorMid := brn.MakeBigFloat(0.0)
	leftSectorMid.Add(&left, &midR)
	leftSectorMid.Quo(&leftSectorMid, &bigTwo)

	rightSectorMid := brn.MakeBigFloat(0.0)
	rightSectorMid.Add(&right, &midR)
	rightSectorMid.Quo(&rightSectorMid, &bigTwo)

	topSectorMid := brn.MakeBigFloat(0.0)
	topSectorMid.Add(&top, &midI)
	topSectorMid.Quo(&topSectorMid, &bigTwo)

	bottomSectorMid := brn.MakeBigFloat(0.0)
	bottomSectorMid.Add(&bottom, &midI)
	bottomSectorMid.Quo(&bottomSectorMid, &bigTwo)

	topSideMid := brn.MakeMember(&bigbase.BigComplex{midR, top})
	bottomSideMid := brn.MakeMember(&bigbase.BigComplex{midR, bottom})
	leftSideMid := brn.MakeMember(&bigbase.BigComplex{left, midI})
	rightSideMid := brn.MakeMember(&bigbase.BigComplex{right, midI})

	topLeftMid := brn.MakeMember(&bigbase.BigComplex{leftSectorMid, topSectorMid})
	topRightMid := brn.MakeMember(&bigbase.BigComplex{rightSectorMid, topSectorMid})
	bottomLeftMid := brn.MakeMember(&bigbase.BigComplex{leftSectorMid, bottomSectorMid})
	bottomRightMid := brn.MakeMember(&bigbase.BigComplex{rightSectorMid, bottomSectorMid})

	tl := bigRegion{
		topLeft:     r.topLeft,
		topRight:    topSideMid,
		bottomLeft:  leftSideMid,
		bottomRight: r.midPoint,
		midPoint:    topLeftMid,
	}
	tr := bigRegion{
		topLeft:     topSideMid,
		topRight:    r.topRight,
		bottomLeft:  r.midPoint,
		bottomRight: rightSideMid,
		midPoint:    topRightMid,
	}
	bl := bigRegion{
		topLeft:     leftSideMid,
		topRight:    r.midPoint,
		bottomLeft:  r.bottomLeft,
		bottomRight: bottomSideMid,
		midPoint:    bottomLeftMid,
	}
	br := bigRegion{
		topLeft:     r.midPoint,
		topRight:    rightSideMid,
		bottomLeft:  bottomSideMid,
		bottomRight: r.bottomRight,
		midPoint:    bottomRightMid,
	}

	brn.subregion = bigSubregion{
		populated: true,
		children:  []bigRegion{tl, tr, bl, br},
	}
}

// RegionMember returns one EscapeValue from the Region
func (brn *BigRegionNumerics) RegionMember() base.EscapeValue {
	return brn.Region.topLeft.EscapeValue
}

func (brn *BigRegionNumerics) proxyNumerics(region *bigRegion) region.RegionNumerics {
	return BigRegionNumericsProxy{
		BigRegionNumerics: brn,
		LocalRegion:   *region,
	}
}

func (brn *BigRegionNumerics) Rect() image.Rectangle {
	return brn.Region.rect(&brn.BigBaseNumerics)
}

func (brn *BigRegionNumerics) SampleDivs() (<-chan uint8, chan<- bool) {
	done := make(chan bool, 1)
	idivch := make(chan uint8, 1)

	go brn.sample(idivch, done)

	return idivch, done
}

func (brn *BigRegionNumerics) sample(idivch chan<- uint8, done <-chan bool) {
	complete := func (idiv uint8) bool {
		select {
		case <-done:
			close(idivch)
			return true
		default:
			idivch<- idiv
			return false
		}
	}

	eval := func (r, i *big.Float) uint8 {
		p := brn.MakeMember(&bigbase.BigComplex{*r, *i})
		p.Mandelbrot(brn.IterateLimit)
		return p.InvDiv
	}

	// Provide the samples we already have
	for _, p := range brn.Points() {
		if complete(p.InvDiv) {
			return
		}
	}

	// Generate samples
	tl := brn.Region.topLeft.C
	br := brn.Region.bottomRight.C
	count := brn.Samples
	fCount := brn.MakeBigFloat(float64(count))
	rmin := tl.Real()
	rmax := br.Real()
	imin := br.Imag()
	imax := tl.Imag()
	width := brn.MakeBigFloat(0.0)
	width.Sub(rmax, rmin)
	height := brn.MakeBigFloat(0.0)
	height.Sub(imax, imin)
	rUnit := brn.MakeBigFloat(0.0)
	rUnit.Quo(&width, &fCount)
	iUnit := brn.MakeBigFloat(0.0)
	iUnit.Quo(&height, &fCount)
	rdown := brn.MakeBigFloat(0.0)
	rdown.Copy(rmax)
	idown := brn.MakeBigFloat(0.0)
	idown.Copy(imax)
	for i := uint(0); i < count; i++ {
		rdown.Sub(&rdown, &rUnit)
		for j := uint(0); j < count; j++ {
			idown.Sub(&idown, &iUnit)
			if complete(eval(&rdown, &idown)) {
				return
			}
		}
	}
	close(idivch)
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

	coords := []bigbase.BigComplex{
		bigbase.BigComplex{left, top},
		bigbase.BigComplex{right, top},
		bigbase.BigComplex{left, bottom},
		bigbase.BigComplex{right, bottom},
		bigbase.BigComplex{midR, midI}, // midpoint is not technically a corner
	}

	points := make([]bigbase.BigEscapeValue, len(coords))

	for i, c := range coords {
		p := bigbase.BigEscapeValue{}
		z := bigbase.MakeBigComplex(0.0, 0.0, prec)
		z.R.Copy(c.Real())
		z.I.Copy(c.Imag())
		p.C = &z
		points[i] = p
	}

	return bigRegion{
		topLeft: points[0],
		topRight: points[1],
		bottomLeft: points[2],
		bottomRight: points[3],
		midPoint: points[4],
	}
}