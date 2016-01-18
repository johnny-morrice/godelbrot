package region

import (
	"fmt"
	"image"
	"functorama.com/demo/base"
	"functorama.com/demo/sequence"
)

type RegionType uint

const (
	UniformPath = RegionType(iota)
	CollapsePath
	SubdividePath
	GlitchPath
)

type MockNumerics struct {
	TExtrinsically 			  bool
	TClaimExtrinsics          bool
	TRect                     bool
	TEvaluateAllPoints        bool
	TSplit                    bool
	TOnGlitchCurve            bool
	TMandelbrotPoints         bool
	TRegionMember             bool
	TSubdivide                bool
	TChildren                 bool
	TRegionSequence bool

	Path                      RegionType

	MockChildren			  []*MockNumerics
	MockSequence			  *MockProxySequence

	AppCollapseSize int
}

func (mock *MockNumerics) Extrinsically(f func()) {
	mock.TExtrinsically = true
	f()
}

func (mock *MockNumerics) ClaimExtrinsics() {
	mock.TClaimExtrinsics = true
}

func (mock *MockNumerics) Rect() image.Rectangle {
	sz := mock.AppCollapseSize
	mock.TRect = true
	if mock.Path == CollapsePath {
		return makeRect(sz)
	} else {
		return makeRect(sz * 10)
	}
}

func makeRect(size int) image.Rectangle {
	return image.Rect(0, 0, size, size)
}

func (mock *MockNumerics) EvaluateAllPoints(iterateLimit uint8) {
	mock.TEvaluateAllPoints = true
}

func (mock *MockNumerics) Split(iterateLimit uint8) {
	mock.TSplit = true
}

func (mock *MockNumerics) OnGlitchCurve(iterateLimit uint8, glitchSamples uint) bool {
	mock.TOnGlitchCurve = true
	return mock.Path == GlitchPath
}

func (mock *MockNumerics) MandelbrotPoints() []base.MandelbrotMember {
	mock.TMandelbrotPoints = true

	const pointCount = 5
	points := make([]base.MandelbrotMember, pointCount)
	switch mock.Path {
	case SubdividePath:
		fallthrough
	case CollapsePath:
		const change = pointCount - 1
		for i := 0; i < change; i++ {
			points[i] = base.BaseMandelbrot{InSet: true, InvDivergence: 0}
		}
		points[change] = base.BaseMandelbrot{InSet: false, InvDivergence: 20}
	case UniformPath:
		fallthrough
	case GlitchPath:
		for i := 0; i < pointCount; i++ {
			points[i] = base.BaseMandelbrot{InSet: true, InvDivergence: 0}
		}
	default:
		panic(fmt.Sprintf("Unknown mock path:", mock.Path))
	}
	return points
}

func (mock *MockNumerics) RegionMember() base.MandelbrotMember {
	mock.TRegionMember = true
	return base.BaseMandelbrot{InSet: true, InvDivergence: 0}
}

func (mock *MockNumerics) Subdivide() bool {
	mock.TSubdivide = true
	return mock.Path == SubdividePath
}

func (mock *MockNumerics) Children() []RegionNumerics {
	mock.TChildren = true
	recurse := make([]RegionNumerics, len(mock.MockChildren))
	for i, child := range mock.MockChildren {
		recurse[i] = child
	}
	return recurse
}

func (mock *MockNumerics) RegionSequence() ProxySequence {
	mock.TRegionSequence = true
	return mock.MockSequence
}

type MockProxySequence struct {
	sequence.MockNumerics

	TClaimExtrinsics bool
	TExtrinsically bool
}

func (mock *MockProxySequence) Extrinsically(f func()) {
	mock.TExtrinsically = true
	f()
}

func (mock *MockProxySequence) ClaimExtrinsics() {
	mock.TClaimExtrinsics = true
}