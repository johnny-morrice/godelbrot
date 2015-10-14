package libgodelbrot

import (
    "image"
)

// This module covers internal details of the numeric strategies used by
// Godelbrot.
// The internal detail of each strategy is defined by an interfac.
// This is because the render control algorithms are also strategies that vary
// independently.

// A full numerics system
type NumericSystem interface {
    RegionNumerics
    SequentialNumerics
}

// Sequential rendering calculations
type SequentialNumerics interface {
    DrawingContext
    OpaqueFlyweightProxy
    MandelbrotSequence()
    ImageDrawSequencer()
    MemberCaptureSequencer()
    CapturedMembers() []PixelMember
}

// Region rendering calculations
type RegionNumerics interface {
    RegionDrawingContext
    OpaqueFlyweightProxy
    Rect() image.Rectange
    EvaluateAllPoints()
    Split()
    Uniform() bool
    OnGlitchCurve() bool
    MandelbrotPoints() []MandelbrotMember
    Subdivide() bool
    Children() []RegionNumerics
    RegionalSequenceNumerics() SequentialNumerics
}

// MandelbrotMember associated with a pixel
type PixelMember struct {
    I int
    J int
    MandelbrotMember
}

// Render the region sequentially
func RenderSequentialRegion(numerics RegionNumerics) {
    smallNumerics := numerics.RegionalSequenceNumerics()
    SequentialRenderImage(smallNumerics)
}

// Subdivide the region, returning true if the subdivision was successful
func Subdivide(numerics RegionNumerics) bool {
    numerics.EvaluateAllPoints()

    if !Uniform(numerics) || numerics.OnGlitchCurve() {
        numerics.Split()
        return true
    }
    return false
}

// Assuming all the region's points() have all been evaluated, 
// return true if they have equal InvDivergence()
func Uniform(numerics RegionNumerics) bool {
    // If inverse divergence on all points is the same, no need to subdivide
    points := numerics.MandelbrotPoints()
    first := points[0].InverseDivergence()
    for _, p := range points[1:] {
        if p.InverseDivergence() != first {
            return false
        }
    }
    return true
}

// Sequence a region, returning points
func SequenceCollapse(numerics RegionNumerics) []PixelMember {
    sequence := region.RegionalSequenceNumerics()
    sequence.ClaimExtrinsics()
    sequence.MemberCaptureSequencer()
    MandelbrotSequence(smallConfig)
    return sequence.CapturedMembers()
}

// Has the region collapsed?
func Collapse(split RegionNumerics) bool {
    rect := split.Rect()
    return rect.Dx() <= split.CollapseSize() || rect.Dy() <= split.CollapseSize()
}