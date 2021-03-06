package region

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/draw"
	"image"
	"testing"
)

func TestNewRegionRenderer(t *testing.T) {
	const iterateLimit uint8 = 200
	const collapse uint = 40
	expectedPic := image.NewNRGBA(image.ZR)
	context := &draw.MockDrawingContext{
		Pic: expectedPic,
	}
	factory := &MockFactory{}
	baseConfig := base.BaseConfig{IterateLimit: iterateLimit}
	regionConfig := RegionConfig{CollapseSize: collapse}
	expectedRenderer := RegionRenderStrategy{
		factory:      factory,
		context:      context,
		regionConfig: regionConfig,
	}
	mock := &MockRenderApplication{}
	mock.RegionFactory = factory
	mock.RegConfig = regionConfig
	mock.Context = context
	mock.Base = baseConfig
	actualRenderer := Make(mock)

	if *actualRenderer != expectedRenderer {
		t.Error("Expected renderer", expectedRenderer,
			"not equal to actual renderer", actualRenderer)
	}

	mockOkay := mock.TRegionNumericsFactory && mock.TDrawingContext
	mockOkay = mockOkay && mock.TRegionConfig

	if !mockOkay {
		t.Error("Expected methods not called on mock")
	}
}

func TestRender(t *testing.T) {
	const iterateLimit uint8 = 200
	const collapseSize int = 40
	expectedPic := image.NewNRGBA(image.ZR)
	mockPalette := &draw.MockPalette{}
	context := &draw.MockDrawingContext{
		Pic: expectedPic,
		Col: mockPalette,
	}
	collapseSequence := &MockProxySequence{}
	uniform := &MockNumerics{
		Path:            UniformPath,
		AppCollapseSize: collapseSize,
	}
	collapse := &MockNumerics{
		Path:            CollapsePath,
		MockSequence:    collapseSequence,
		AppCollapseSize: collapseSize,
	}
	mockNumerics := &MockNumerics{
		Path:            SubdividePath,
		MockChildren:    []*MockNumerics{uniform, collapse},
		MockSequence:    &MockProxySequence{},
		AppCollapseSize: collapseSize,
	}
	factory := &MockFactory{Numerics: mockNumerics}
	regionConfig := RegionConfig{CollapseSize: uint(collapseSize)}
	renderer := RegionRenderStrategy{
		factory:      factory,
		context:      context,
		regionConfig: regionConfig,
	}

	actualPic, err := renderer.Render()

	if actualPic != expectedPic {
		t.Error("Expected pic differed from actual:", actualPic)
	}

	if err != nil {
		t.Error("Unexpeced error in render:", err)
	}

	if !factory.TBuild {
		t.Error("Expected methods not called on factory:", factory)
	}

	if !(uniform.TClaimExtrinsics && uniform.TRect) {
		t.Error("Expected methods not called on uniform region:", uniform)
	}

	if !(collapse.TClaimExtrinsics && collapse.TRegionSequence) {
		t.Error("Expected methods not called on collapse region:", collapse)
	}

	if !mockPalette.TColor {
		t.Error("Expected methods not called on paleete:", mockPalette)
	}

	sequenceOkay := collapseSequence.TExtrinsically && collapseSequence.TSequence
	if !sequenceOkay {
		t.Error("Expected methods not called on collapsed sequence numerics:", collapseSequence)
	}
}

func TestSubdivideRegions(t *testing.T) {
	const iterateLimit uint8 = 200
	const collapseSize = 40
	uniform := newMockNumerics(UniformPath, collapseSize)
	collapse := newMockNumerics(CollapsePath, collapseSize)
	mock := &MockNumerics{
		Path:            SubdividePath,
		MockChildren:    []*MockNumerics{uniform, collapse},
		AppCollapseSize: collapseSize,
	}
	regionConfig := RegionConfig{CollapseSize: collapseSize}
	renderer := RegionRenderStrategy{
		regionConfig: regionConfig,
	}

	actualUniform, actualCollapse := renderer.SubdivideRegions(mock)

	actualUniCount := len(actualUniform)
	const expectUniCount = 1
	if actualUniCount != expectUniCount {
		t.Error("Expected ", expectUniCount, "uniform regions but received", actualUniCount)
	}

	actualCollCount := len(actualCollapse)
	const expectCollCount = 1
	if actualUniCount != expectCollCount {
		t.Error("Expected ", expectUniCount, "collapsed regions but received", actualCollCount)
	}

	mockOkay := mock.TClaimExtrinsics && mock.TRect
	mockOkay = mockOkay && mock.TMandelbrotPoints && mock.TSplit
	mockOkay = mockOkay && mock.TChildren
	if !mockOkay {
		t.Error("Expected methods not called on inital region:", mock)
	}

	uniformOkay := uniform.TClaimExtrinsics && uniform.TRect
	uniformOkay = uniformOkay && uniform.TMandelbrotPoints
	if !uniformOkay {
		t.Error("Expected methods not called on uniform region:", uniform)
	}

	collapseOkay := collapse.TClaimExtrinsics && collapse.TRect
	if !collapseOkay {
		t.Error("Expected methods not called on collapsed region:", collapse)
	}
}
