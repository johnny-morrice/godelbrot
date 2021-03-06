package nativeregion

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/nativebase"
	"testing"
)

func TestRegionSplitPos(t *testing.T) {
	helper := NativeRegionSplitHelper{
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
	helper := NativeRegionSplitHelper{
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
	helper := NativeRegionSplitHelper{
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
	const inputChildCount = 4
	numerics := NativeRegionNumerics{
		subregion: nativeSubregion{
			populated: true,
			// We are not inspecting the children here
			children: make([]nativeRegion, inputChildCount),
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

	triggerPanic()

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

func TestRect(t *testing.T) {
	const picSide = 2
	const planeSide = 2.0

	const left = -1
	const bottom = -1
	const right = left + planeSide
	const top = bottom + planeSide

	min := complex(left, bottom)
	max := complex(right, top)
	parent := nativebase.NativeBaseNumerics{
		BaseNumerics: base.BaseNumerics{
			PicXMin: 0,
			PicXMax: picSide,
			PicYMin: 0,
			PicYMax: picSide,
		},
		SqrtDivergeLimit: sqrtDLimit,
		RealMin:          left,
		RealMax:          right,
		ImagMin:          bottom,
		ImagMax:          top,
	}
	numerics := &NativeRegionNumerics{
		NativeBaseNumerics: parent,
	}
	numerics.Region = createNativeRegion(parent, min, max)

	uq := nativebase.UnitQuery{picSide, picSide, planeSide, planeSide}
	numerics.Runit, numerics.Iunit = uq.PixelUnits()

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

func testRegionSplit(helper NativeRegionSplitHelper, t *testing.T) {
	const iterlim = uint8(255)
	parent := nativebase.NativeBaseNumerics{}
	parent.SqrtDivergeLimit = sqrtDLimit
	parent.IterateLimit = iterlim

	initMin := complex(helper.left, helper.bottom)
	initMax := complex(helper.right, helper.top)

	topLeftMin := complex(helper.left, helper.midI)
	topLeftMax := complex(helper.midR, helper.top)

	topRightMin := complex(helper.midR, helper.midI)
	topRightMax := complex(helper.right, helper.top)

	bottomLeftMin := complex(helper.left, helper.bottom)
	bottomLeftMax := complex(helper.midR, helper.midI)

	bottomRightMin := complex(helper.midR, helper.bottom)
	bottomRightMax := complex(helper.right, helper.midI)

	subjectRegion := createNativeRegion(parent, initMin, initMax)

	expected := []nativeRegion{
		createNativeRegion(parent, topLeftMin, topLeftMax),
		createNativeRegion(parent, topRightMin, topRightMax),
		createNativeRegion(parent, bottomLeftMin, bottomLeftMax),
		createNativeRegion(parent, bottomRightMin, bottomRightMax),
	}

	numerics := NativeRegionNumerics{
		Region: subjectRegion,
	}
	numerics.NativeBaseNumerics = parent
	numerics.Split()
	actualChildren := numerics.subregion.children

	for i, ex := range expected {
		actual := actualChildren[i]
		exPoints := ex.points()
		acPoints := actual.points()
		fail := false
		for j, e := range exPoints {
			a := acPoints[j]
			if e.C != a.C {
				fail = true
				t.Log("Region", i, "error at point", j,
					"expected", e,
					"but received", a)
			}
		}
		if fail {
			t.Fail()
		}
	}
}

type NativeRegionSplitHelper struct {
	left   float64
	right  float64
	bottom float64
	top    float64
	midR   float64
	midI   float64
}
