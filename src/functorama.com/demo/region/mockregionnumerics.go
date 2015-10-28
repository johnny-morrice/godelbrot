package region

import (
	"image"
)

type regionType uint

const (
	uniform = regionType(iota)
	collapse
	subdivide
	glitch
)

type mockRegionNumerics struct {
	path                      regionType
	tClaimExtrinsics          bool
	tRect                     bool
	tEvaluateAllPoints        bool
	tSplit                    bool
	tOnGlitchCurve            bool
	tMandelbrotPoints         bool
	tRegionMember             bool
	tSubdivide                bool
	tChildren                 bool
	tRegionalSequenceNumerics bool
}

const collapseSize = 20
const collapseMembers = collapseSize * collapseSize // C.M.
const children = 4                                  // C

func (mock mockRegionNumerics) ClaimExtrinsics() {
	mock.tClaimExtrinsics = true
}

func (mock mockRegionNumerics) Rect() image.Rectangle {
	mock.tRect = true
	if mock.path == collapse {
		return makeRect(collapseSize)
	} else {
		return makeRect(collapseMembers)
	}
}

func (mock mockRegionNumerics) EvaluateAllPoints(iterateLimit int) {
	mock.tEvaluateAllPoints = true
}

func (mock mockRegionNumerics) Split() {
	mock.tSplit = true
}

func (mock mockRegionNumerics) OnGlitchCurve(iterateLimit uint8, glitchSamples uint) bool {
	mock.tOnGlitchCurve = true
	return mock.path == glitch
}

type pointFunc func(mockRegionNumerics) MandelbrotMember

func (mock mockRegionNumerics) MandelbrotPoints() []MandelbrotMember {
	mock.tMandelbrotPoints = true
	fs := []pointFunc{
		mockTopLeft,
		mockTopRight,
		mockBottomLeft,
		mockBottomRight,
		mockMidPoint,
	}

	points := make([]MandelbrotMember, len(fs))

	for i, f := range fs {
		points[i] = f(mock)
	}

	return points
}

func (mock mockRegionNumerics) RegionMember() {
	mock.tRegionMember = true
	return mockMidPoint(mock)
}

func (mock mockRegionNumerics) Subdivide() bool {
	mock.tSubdivide = true
	return mock.path == subdivide
}

func (mock mockRegionNumerics) Children() []RegionNumerics {
	mock.tChildren = true
	recurse := make([]RegionNumerics, children)
	for i := 0; i < children; i++ {
		recurse[i] = mock
	}
	return recurse
}

func (mock mockRegionNumerics) RegionalSequenceNumerics() SequentialNumerics {
	mock.tRegionalSequenceNumerics = true
	mockSequence := mockSequenceNumerics{
		minR: 0,
		maxR: collapseSize,
		minI: 0,
		maxI: collapseSize,
	}
	return mockSequence
}
