package region

import (
	"github.com/johnny-morrice/godelbrot/internal/draw"
	"testing"
)

func TestRenderSequenceRegion(t *testing.T) {
	const iterateLimit = 40
	mockSequence := &MockProxySequence{}
	mockRegion := &MockNumerics{
		MockSequence: mockSequence,
	}
	// Not inspecting contract re DrawingContext
	context := draw.NewMockDrawingContext(iterateLimit)
	RenderSequenceRegion(mockRegion, context)

	if !mockRegion.TRegionSequence {
		t.Error("Expected methods not called on mockRegion:", mockRegion)
	}

	sequenceOkay := mockSequence.TExtrinsically && mockSequence.TSequence
	if !sequenceOkay {
		t.Error("Expected methods not called on mockSequence:", mockSequence)
	}
}

func TestSequenceCollapse(t *testing.T) {
	const iterateLimit = 40

	mockSequence := &MockProxySequence{}
	mockRegion := &MockNumerics{
		MockSequence: mockSequence,
	}
	SequenceCollapse(mockRegion)
	if !mockRegion.TRegionSequence {
		t.Error("Expected methods not called on mockRegion:", mockRegion)
	}

	sequenceOkay := mockSequence.TExtrinsically && mockSequence.TSequence
	if !sequenceOkay {
		t.Error("Expected methods not called on mockSequence:", mockSequence)
	}
}

func TestSubdivide(t *testing.T) {
	const iterateLimit uint8 = 200
	const glitchSamples uint = 10
	uniReg := &MockNumerics{Path: UniformPath}
	// Collapsed regions shouldn't care about subdivision
	//collReg := MockNumerics{path: collapse}
	subReg := &MockNumerics{Path: SubdividePath}

	const (
		uniIndex = iota
		// collIndex
		subIndex
		gliIndex
	)
	regions := []RegionNumerics{uniReg, subReg}
	actual := make([]bool, len(regions))

	for i, reg := range regions {
		actual[i] = Subdivide(reg)
	}

	// Results for uniform region
	if actual[uniIndex] {
		t.Error("Expected negative Subdivide return for uniform region")
	}
	if uniReg.TSplit {
		t.Error("Split was called on uniform region")
	}
	// We expect the mandelbrot points to be examined for uniformity detection
	if !(uniReg.TMandelbrotPoints) {
		t.Error("Expected methods were not called on uniform region:", uniReg)
	}

	// Results for subdividing region
	if !actual[subIndex] {
		t.Error("Expected positive Subdivide return for subdivided region")
	}
	if !subReg.TSplit {
		t.Error("Expected methods were not called on subdivided region:", subReg)
	}
}

func TestUniform(t *testing.T) {
	uniform := &MockNumerics{Path: UniformPath}
	collapse := &MockNumerics{Path: CollapsePath}
	subdivide := &MockNumerics{Path: SubdividePath}

	numerics := []RegionNumerics{uniform, collapse, subdivide}
	expect := []bool{true, false, false}

	for i, num := range numerics {
		actual := Uniform(num)
		exp := expect[i]
		if exp != actual {
			if exp {
				t.Error("Expected numeric", i, "to have positive Uniform return")
			} else {
				t.Error("Expected numeric", i, "to have negative Uniform return")
			}
		}
	}
}

func newMockNumerics(path RegionType, collapseSize int) *MockNumerics {
	return &MockNumerics{
		Path:            path,
		AppCollapseSize: collapseSize,
	}
}

func TestCollapse(t *testing.T) {
	const collapseSize = 10
	uniform := newMockNumerics(UniformPath, collapseSize)
	collapse := newMockNumerics(CollapsePath, collapseSize)
	subdivide := newMockNumerics(SubdividePath, collapseSize)

	numerics := []RegionNumerics{uniform, subdivide, collapse}
	expect := []bool{false, false, true}

	for i, num := range numerics {
		actual := Collapse(num, collapseSize)
		exp := expect[i]
		if exp != actual {
			if exp {
				t.Error("Expected numeric", i, "to have positive Collapse return")
			} else {
				t.Error("Expected numeric", i, "to have negative Collapse return")
			}
		}
	}
}
