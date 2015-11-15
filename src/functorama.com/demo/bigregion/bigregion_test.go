package bigregion

import (
	"math/big"
	"testing"
)

func TestRegionSplitPos(t *testing.T) {
	helper := bigRegionSplitHelper{
		left:   CreateBigFloat(1.0, Prec64),
		right:  CreateBigFloat(3.0, Prec64),
		top:    CreateBigFloat(3.0, Prec64),
		bottom: CreateBigFloat(1.0, Prec64),
		midR:   CreateBigFloat(2.0, Prec64),
		midI:   CreateBigFloat(2.0, Prec64),
	}

	testRegionSplit(helper, t)
}

func TestRegionSplitNeg(t *testing.T) {
	helper := bigRegionSplitHelper{
		left:   CreateBigFloat(-100.0, Prec64),
		right:  CreateBigFloat(-24.0, Prec64),
		top:    CreateBigFloat(-10.0, Prec64),
		bottom: CreateBigFloat(-340.0, Prec64),
		midR:   CreateBigFloat(-62.0, Prec64),
		midI:   CreateBigFloat(-175.0, Prec64),
	}

	testRegionSplit(helper, t)
}

func TestRegionSplitNegPos(t *testing.T) {
	helper := bigRegionSplitHelper{
		left:   CreateBigFloat(-100.0, Prec64),
		right:  CreateBigFloat(24.0, Prec64),
		top:    CreateBigFloat(10.0, Prec64),
		bottom: CreateBigFloat(-340.0, Prec64),
		midR:   CreateBigFloat(-38.0, Prec64),
		midI:   CreateBigFloat(-165.0, Prec64),
	}

	testRegionSplit(helper, t)
}

func TestChildrenPopulated(t *testing.T) {
	numerics := BigRegionNumerics{
		subRegion: bigSubregion{
			populated: true,
			// We are not inspecting the children here
			children: []bigRegion{nil, nil, nil, nil},
		},
	}
	children := numerics.Children()

	expectedChildren := 4
	actualChildren := len(children)
	if actualChildren != expectedChildren {
		t.Error("Expected", expectedChildren, "but received", actualChildren)
	}
}

func TestChildrenEmpty(t *testing.T) {
	numerics := BigRegionNumerics{}

	recovered := false
	triggerPanic := func() {
		defer func() {
			r := recover()
			recovered = r != nil
		}()
		numerics.Children()
	}

	if !recovered {
		t.Error("Expected panic e.g. \"Error when raising error\"")
	}
}

func TestMandelbrotPoints(t *testing.T) {
	numerics := BigRegionNumerics{}

	points := numerics.MandelbrotPoints()

	expectedPoints := 5
	actualPoints := len(points)

	if expectedPoints != actualPoints {
		t.Error("Expected to receive", expectedPoints, "but received", actualPoints)
	}
}

func TestEvaluateAllPoints(t *testing.T) {
	numerics := BigRegionNumerics{}
	numerics.EvaluateAllPoints(1)
	region := numerics.region

	okay := region.topLeft.evaluated
	okay = okay && region.topRight.evaluated
	okay = okay && region.bottomLeft.evaluated
	okay = okay && region.bottomRight.evaluated
	okay = okay && region.midPoint.evaluated

	if !okay {
		t.Error("Expected all points to be evaluated, but region was:", region)
	}
}

func TestRect(t *testing.T) {
	left := -1.0
	bottom := -1.0
	right := 1.0
	top := 1.0

	min := CreateBigComplex(left, bottom)
	max := CreateBigComplex(right, top)
	numerics := BigRegionNumerics{
		region:  createBigRegion(min, max),
		picXMin: 0,
		picXMax: 2,
		picYMin: 0,
		picYMax: 2,
		realMin: CreateBigFloat(left, Prec64),
		realMax: CreateBigFloat(right, Prec64),
		imagMin: CreateBigFloat(bottom, Prec64),
		imagMax: CreateBigFloat(top, Prec64),
	}

	expectMinX := 0
	expectMaxX := 2
	expectMinY := 0
	expectMaxY := 2

	r := numerics.Rect()

	okay := expectMinX == r.Min.X
	okay = okay && expectMaxX == r.Max.X
	okay = okay && expectMinY == r.Min.Y
	okay = okay && expectMaxY == r.Max.Y
	if !okay {
		t.Error("Rectangle had unexpected bounds", r)
	}

}

func testRegionSplit(helper bigRegionSplitHelper, t *testing.T) {
	initMin := BigComplex{helper.left, helper.bottom}
	initMax := BigComplex{helper.right, helper.top}

	topLeftMin := BigComplex{helper.left, midI}
	topLeftMax := BigComplex{midR, helper.top}

	topRightMin := BigComplex{midR, midI}
	topRightMax := BigComplex{helper.right, helper.top}

	bottomLeftMin := BigComplex{helper.left, helper.bottom}
	bottomLeftMax := BigComplex{midR, midI}

	bottomRightMin := BigComplex{midR, helper.bottom}
	bottomRightMax := BigComplex{helper.right, midI}

	subjectRegion := createBigRegion(initMin, initMax)

	expected := []Region{
		createBigRegion(topLeftMin, topLeftMax),
		createBigRegion(topRightMin, topRightMax),
		createBigRegion(bottomLeftMin, bottomLeftMax),
		createBigRegion(bottomRightMin, bottomRightMax),
	}

	actual := subjectRegion.Split()

	for i, ex := range expected {
		similarity := sameRegion(ex, actual.children[i])
		if !similarity.same {
			t.Error(
				"Unexpected child region ", i,
				", expected point ", similarity.n,
				"to be ", similarity.a,
				" but was ", similarity.b,
			)
		}
	}
}

type bigRegionSplitHelper struct {
	left   big.Float
	right  big.Float
	bottom big.Float
	top    big.Float
	midR   big.Float
	midI   big.Float
}

type bigRegionSameness struct {
	a            big.Float
	b            complex128
	regionNumber int
	same         bool
}

func sameRegion(a Region, b Region) bigRegionSameness {
	aPoints := bigPoints(b)
	bPoints := bigPoints(a)

	for i, ap := range aPoints {
		bp := bPoints[i]
		if ap.C.Cmp(bp.C) {
			return bigRegionSameness{
				a:            ap.c,
				b:            bp.c,
				regionNumber: i,
				same:         false,
			}
		}
	}
	return bigRegionSameness{same: true}
}
