package base

import (
    "testing"
)

func TestRenderSequentialRegion(t *testing.T) {
    mockRegion := mockRegionNumerics{}
    mockSequence := mockRegion.RegionalSequenceNumerics()
    // What methods do we expect in a sequence render?
    // Sequence rendering in full is beyond the contract of this method
    if !mockSequence.tClaimExtrinsics {
        t.Error("Expected mockSequence to ClaimExtrinsics")
    }
}

func TestSequenceCollapse(t *testing.T) {

}

func TestSubdivide(t *testing.T) {
    uniReg := mockRegionNumerics{path: uniform}
    // Collapsed regions shouldn't care about subdivision
    //collReg := mockRegionNumerics{path: collapse}
    subReg := mockRegionNumerics{path: subdivide}
    gliReg := mockRegionNumerics{path: glitch}

    const (
        uniIndex = iota
        // collIndex
        subIndex
        gliIndex
    )
    regions := []RegionNumerics{uniReg, collReg, subReg, gliReg}
    actual := make([]bool, len(regions))

    for i, reg := range regions {
        actual[i] = Subdivide(reg)
    }

    // Results for uniform region
    if actual[uniIndex] {
        t.Error("Expected negative Subdivide return for uniform region")
    }
    if uniReg.tSplit {
        t.Error("Split was called on uniform region")
    }
    // We expect the mandelbrot points to be examined for uniformity detection
    if !(uniReg.tEvaluateAllPoints && uniReg.tMandelbrotPoints) {
        t.Error("Expected methods were not called on uniform region:", uniReg)
    }

    // Results for glitch region
    if actual[gliIndex] {
        t.Error("Expected negative Subdivide return for glitched region")
    }
    if !(gliReg.tEvaluateAllPoints && gliReg.tOnGlitchCurve) {
        t.Error("Expected methods were not called on glitched region")
    }

    // Results for subdividing region
    if !actual[subIndex] {
        t.Error("Expected positive Subdivide return for subdivided region")
    }
    if !(subReg.tSplit && subReg.tEvaluateAllPoints && subReg.tUniform && subReg.tOnGlitchCurve) {
        t.Error("Expected method were not called on subdivided region:", subReg)
    }
}
