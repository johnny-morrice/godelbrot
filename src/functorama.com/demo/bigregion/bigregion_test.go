package bigregion

import (
	"math/big"
	"testing"
	"functorama.com/demo/bigbase"
)

const prec = 53

func TestRegionSplitPos(t *testing.T) {
	helper := bigRegionSplitHelper{
		left:   bigbase.MakeBigFloat(1.0, prec),
		right:  bigbase.MakeBigFloat(3.0, prec),
		top:    bigbase.MakeBigFloat(3.0, prec),
		bottom: bigbase.MakeBigFloat(1.0, prec),
		midR:   bigbase.MakeBigFloat(2.0, prec),
		midI:   bigbase.MakeBigFloat(2.0, prec),
	}

	testRegionSplit(helper, t)
}

func TestRegionSplitNeg(t *testing.T) {
	helper := bigRegionSplitHelper{
		left:   bigbase.MakeBigFloat(-100.0, prec),
		right:  bigbase.MakeBigFloat(-24.0, prec),
		top:    bigbase.MakeBigFloat(-10.0, prec),
		bottom: bigbase.MakeBigFloat(-340.0, prec),
		midR:   bigbase.MakeBigFloat(-62.0, prec),
		midI:   bigbase.MakeBigFloat(-175.0, prec),
	}

	testRegionSplit(helper, t)
}

func TestRegionSplitNegPos(t *testing.T) {
	helper := bigRegionSplitHelper{
		left:   bigbase.MakeBigFloat(-100.0, prec),
		right:  bigbase.MakeBigFloat(24.0, prec),
		top:    bigbase.MakeBigFloat(10.0, prec),
		bottom: bigbase.MakeBigFloat(-340.0, prec),
		midR:   bigbase.MakeBigFloat(-38.0, prec),
		midI:   bigbase.MakeBigFloat(-165.0, prec),
	}

	testRegionSplit(helper, t)
}

func TestChildrenPopulated(t *testing.T) {
	const childCount = 4

	numerics := BigRegionNumerics{
		subregion: bigSubregion{
			populated: true,
			// We are not inspecting the children here
			children: make([]bigRegion, childCount),
		},
	}
	children := numerics.Children()

	actualChildren := len(children)
	if actualChildren != childCount {
		t.Error("Expected", childCount, "but received", actualChildren)
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

	triggerPanic()

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
	const iterLimit = 1
	numerics := BigRegionNumerics{}
	numerics.SqrtDivergeLimit = bigbase.MakeBigFloat(2.0, prec)

	for _, thunk := range numerics.Region.thunks() {
		thunk.C = &thunk.cStore
		thunk.SqrtDivergeLimit = &numerics.SqrtDivergeLimit
	}

	numerics.EvaluateAllPoints(iterLimit)
	region := numerics.Region

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

	app := &MockRenderApplication{}
	app.UserMin = bigbase.MakeBigComplex(left, bottom, prec)
	app.UserMax = bigbase.MakeBigComplex(right, top, prec)
	app.Prec = 53
	app.PictureWidth = 2
	app.PictureHeight = 2

	numerics := Make(app)

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
	initMin := bigbase.BigComplex{helper.left, helper.bottom}
	initMax := bigbase.BigComplex{helper.right, helper.top}

	topLeftMin := bigbase.BigComplex{helper.left, helper.midI}
	topLeftMax := bigbase.BigComplex{helper.midR, helper.top}

	topRightMin := bigbase.BigComplex{helper.midR, helper.midI}
	topRightMax := bigbase.BigComplex{helper.right, helper.top}

	bottomLeftMin := bigbase.BigComplex{helper.left, helper.bottom}
	bottomLeftMax := bigbase.BigComplex{helper.midR, helper.midI}

	bottomRightMin := bigbase.BigComplex{helper.midR, helper.bottom}
	bottomRightMax := bigbase.BigComplex{helper.right, helper.midI}

	subjectRegion := createBigRegion(initMin, initMax)

	t.Log("Subject middle:", bigbase.DbgC(subjectRegion.midPoint.cStore))

	expected := []bigRegion{
		createBigRegion(topLeftMin, topLeftMax),
		createBigRegion(topRightMin, topRightMax),
		createBigRegion(bottomLeftMin, bottomLeftMax),
		createBigRegion(bottomRightMin, bottomRightMax),
	}

	numerics := BigRegionNumerics{}
	numerics.Region = subjectRegion
	numerics.SqrtDivergeLimit = bigbase.MakeBigFloat(2.0, prec)
	numerics.Precision = prec
	numerics.Split()
	actualChildren := numerics.subregion.children

	for i, ex := range expected {
		actual := actualChildren[i]
		exPoints := ex.thunks()
		acPoints := actual.thunks()

		fail := false
		for j, expectThunk := range exPoints {
			actThunk := acPoints[j]
			okay := bigbase.BigComplexEq(&expectThunk.cStore, &actThunk.cStore)
			okay = okay && bigbase.BigComplexEq(expectThunk.C, actThunk.C)
			if !okay {
				fail = true
				t.Log("Region", i, "error at thunk", j,
					"\nexpected\t", bigbase.DbgC(expectThunk.cStore),
					"\nbut received\t", bigbase.DbgC(actThunk.cStore))
			}
		}
		if fail {
			t.Fail()
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