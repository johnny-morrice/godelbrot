package region

import (
	"fmt"
	"image"
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/sequence"
)

type RegionType uint

const (
	UniformPath = RegionType(iota)
	CollapsePath
	SubdividePath
)

type MockNumerics struct {
	TExtrinsically 			  bool
	TClaimExtrinsics          bool
	TRect                     bool
	TSplit                    bool
	TMandelbrotPoints         bool
	TRegionMember             bool
	TSubdivide                bool
	TChildren                 bool
	TSampleDivs				  bool
	TRegionSequence bool

	Path                      RegionType

	MockChildren			  []*MockNumerics
	MockSequence			  *MockProxySequence

	AppCollapseSize int
}

func (mock *MockNumerics) SampleDivs() (<-chan uint8, chan<- bool) {
	mock.TSampleDivs = true
	done := make(chan bool, 1)
	idivch := make(chan uint8, 1)

	go func() {
		for _, p := range mock.MandelbrotPoints() {
			idivch<- p.InvDiv
		}
		close(idivch)
	}()

	return idivch, done
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

func (mock *MockNumerics) Split() {
	mock.TSplit = true
}

func (mock *MockNumerics) MandelbrotPoints() []base.EscapeValue {
	mock.TMandelbrotPoints = true

	const pointCount = 5
	points := make([]base.EscapeValue, pointCount)
	switch mock.Path {
	case SubdividePath:
		fallthrough
	case CollapsePath:
		const change = pointCount - 1
		for i := 0; i < change; i++ {
			points[i] = base.EscapeValue{InSet: true, InvDiv: 0}
		}
		points[change] = base.EscapeValue{InSet: false, InvDiv: 20}
	case UniformPath:
		for i := 0; i < pointCount; i++ {
			points[i] = base.EscapeValue{InSet: true, InvDiv: 0}
		}
	default:
		panic(fmt.Sprintf("Unknown mock path:", mock.Path))
	}
	return points
}

func (mock *MockNumerics) RegionMember() base.EscapeValue {
	mock.TRegionMember = true
	return base.EscapeValue{InSet: true, InvDiv: 0}
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