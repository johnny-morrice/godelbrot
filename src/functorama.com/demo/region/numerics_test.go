package region

import (
    "testing"
    "functorama.com/demo/base"
)

func TestRenderSequenceRegion(t *testing.T) {
    const iterateLimit = 40
    mockSequence := &MockProxySequence{}
    mockRegion := &MockNumerics{
        MockSequence: mockSequence,
    }
    // Not inspecting contract re DrawingContext
    RenderSequenceRegion(mockRegion, nil, iterateLimit)

    if !mockRegion.TRegionSequence {
        t.Error("Expected methods not called on mockRegion:", mockRegion)
    }

    sequenceOkay := mockSequence.TClaimExtrinsics && mockSequence.TImageDrawSequencer
    sequenceOkay = sequenceOkay && mockSequence.TMandelbrotSequence
    if !sequenceOkay {
        t.Error("Expected methods not called on mockSequence:", mockSequence)
    }
}

func TestSequenceCollapse(t *testing.T) {
    const iterateLimit = 40
    expectedMembers := []base.PixelMember{
        base.PixelMember{I: 20, J: 40},
    }
    mockSequence := &MockProxySequence{}
    mockSequence.Captured = expectedMembers
    mockRegion := &MockNumerics{
        MockSequence: mockSequence,
    }
    actualMembers := SequenceCollapse(mockRegion, iterateLimit)

    membersOkay := len(actualMembers) == len(expectedMembers)
    for i, exp := range expectedMembers {
        membersOkay = membersOkay && exp == actualMembers[i]
    }

    if !membersOkay {
        t.Error("Expected members", expectedMembers, "but received:", actualMembers)
    }

    if !mockRegion.TRegionSequence {
        t.Error("Expected methods not called on mockRegion:", mockRegion)
    }

    sequenceOkay := mockSequence.TClaimExtrinsics && mockSequence.TMemberCaptureSequencer
    sequenceOkay = sequenceOkay && mockSequence.TMandelbrotSequence && mockSequence.TCapturedMembers
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
    gliReg := &MockNumerics{Path: GlitchPath}

    const (
        uniIndex = iota
        // collIndex
        subIndex
        gliIndex
    )
    regions := []RegionNumerics{uniReg, subReg, gliReg}
    actual := make([]bool, len(regions))

    for i, reg := range regions {
        actual[i] = Subdivide(reg, iterateLimit, glitchSamples)
    }

    // Results for uniform region
    if actual[uniIndex] {
        t.Error("Expected negative Subdivide return for uniform region")
    }
    if uniReg.TSplit {
        t.Error("Split was called on uniform region")
    }
    // We expect the mandelbrot points to be examined for uniformity detection
    if !(uniReg.TEvaluateAllPoints && uniReg.TMandelbrotPoints && uniReg.TOnGlitchCurve) {
        t.Error("Expected methods were not called on uniform region:", uniReg)
    }

    // Results for glitch region
    if !actual[gliIndex] {
        t.Error("Expected positive Subdivide return for glitched region")
    }
    if !(gliReg.TEvaluateAllPoints && gliReg.TMandelbrotPoints && gliReg.TOnGlitchCurve) {
        t.Error("Expected methods were not called on glitched region")
    }

    // Results for subdividing region
    if !actual[subIndex] {
        t.Error("Expected positive Subdivide return for subdivided region")
    }
    if !(subReg.TSplit && subReg.TEvaluateAllPoints) {
        t.Error("Expected methods were not called on subdivided region:", subReg)
    }
}

func TestUniform(t *testing.T) {
    uniform := &MockNumerics{Path: UniformPath}
    collapse := &MockNumerics{Path: CollapsePath}
    subdivide := &MockNumerics{Path: SubdividePath}
    glitch := &MockNumerics{Path: GlitchPath}

    numerics := []RegionNumerics{uniform, glitch, collapse, subdivide}
    expect := []bool{true, true, false, false}

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
        Path: path,
        AppCollapseSize: collapseSize,
    }
}

func TestCollapse(t *testing.T) {
    const collapseSize = 10
    uniform := newMockNumerics(UniformPath, collapseSize)
    collapse := newMockNumerics(CollapsePath, collapseSize)
    subdivide := newMockNumerics(SubdividePath, collapseSize)
    glitch := newMockNumerics(GlitchPath, collapseSize)

    numerics := []RegionNumerics{uniform, glitch, subdivide, collapse}
    expect := []bool{false, false, false, true}

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
