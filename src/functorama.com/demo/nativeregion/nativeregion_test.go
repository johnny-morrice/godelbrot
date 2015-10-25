package region

import (
	"math/big"
	"testing"
)

func TestRegionSplitPos(t *testing.T) {
	helper := nativeRegionSplitHelper{
		left:   1.0,
		right:  3.0,
		top:    3.0,
		bottom: 1.0,
		midR:   2.0,
		midI:   2.0,
	}

	testRegionSplit(helper, t)
}

func TestRegionSplitNeg(t *testing.T) {
	helper := nativeRegionSplitHelper{
		left:   -100.0,
		right:  -24.0,
		top:    -10.0,
		bottom: -340.0,
		midR:   -62.0,
		midI:   -175.0,
	}

	testRegionSplit(helper, t)
}

func TestRegionSplitNegPos(t *testing.T) {
	helper := nativeRegionSplitHelper{
		left:   -100.0,
		right:  24.0,
		top:    10.0,
		bottom: -340.0,
		midR:   -38.0,
		midI:   -165.0,
	}

	testRegionSplit(helper, t)
}

func TestChildrenPopulated(t *testing.T) {
	numerics := NativeRegionNumerics{
		subRegion: nativeSubregion{
			populated: true,
			// We are not inspecting the children here
			children: []nativeRegion{nil, nil, nil, nil},
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
	numerics := NativeRegionNumerics{}

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
	numerics := NativeRegionNumerics{}

	points := numerics.MandelbrotPoints()

	expectedPoints := 5
	actualPoints := len(points)

	if expectedPoints != actualPoints {
		t.Error("Expected to receive", expectedPoints, "but received", actualPoints)
	}
}

func TestEvaluateAllPoints(t *testing.T) {
	numerics := NativeRegionNumerics{}
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

	min := complex(left, bottom)
	max := complex(right, top)
	numerics := NativeRegionNumerics{
		region:  createNativeRegion(min, max),
		picXMin: 0,
		picXMax: 2,
		picYMin: 0,
		picYMax: 2,
		realMin: left,
		realMax: right,
		imagMin: bottom,
		imagMax: top,
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

func testRegionSplit(helper nativeRegionSplitHelper, t *testing.T) {
	initMin := complex(helper.left, helper.bottom)
	initMax := complex(helper.right, helper.top)

	topLeftMin := complex(helper.left, midI)
	topLeftMax := complex(midR, helper.top)

	topRightMin := complex(midR, midI)
	topRightMax := complex(helper.right, helper.top)

	bottomLeftMin := complex(helper.left, helper.bottom)
	bottomLeftMax := complex(midR, midI)

	bottomRightMin := complex(midR, helper.bottom)
	bottomRightMax := complex(helper.right, midI)

	subjectRegion := createNativeRegion(initMin, initMax)

	expected := []Region{
		createNativeRegion(topLeftMin, topLeftMax),
		createNativeRegion(topRightMin, topRightMax),
		createNativeRegion(bottomLeftMin, bottomLeftMax),
		createNativeRegion(bottomRightMin, bottomRightMax),
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

type nativeRegionSplitHelper struct {
	left   float64
	right  float64
	bottom float64
	top    float64
	midR   float64
	midI   float64
}

type nativeRegonSameness struct {
	a            float64
	b            complex128
	regionNumber int
	same         bool
}

func sameRegion(a Region, b Region) nativeRegionSameness {
	aPoints := nativePoints(b)
	bPoints := nativePoints(a)

	for i, ap := range aPoints {
		bp := bPoints[i]
		if ap.C.Cmp(bp.C) {
			return nativeRegionSameness{
				a:            ap.c,
				b:            bp.c,
				regionNumber: i,
				same:         false,
			}
		}
	}
	return nativeRegionSameness{same: true}
}
