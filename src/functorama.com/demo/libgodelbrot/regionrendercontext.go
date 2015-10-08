package libgodelbrot

import (
    "image"
)

type RegionRenderContext interface {
    Collapse() bool
    Rect() image.Rectange
    EvaluateAllPoints()
    Split()
    Uniform() bool
    OnGlitchCurve() bool
    MandelbrotPoints() []MandelbrotMember
    Subdivide() bool
    Children() []RegionRenderContext
}

func Subdivide(context RegionRenderContext) {
    context.EvaluateAllPoints()

    if !Uniform(context) || context.OnGlitchCurve() {
        context.Split()
        return true
    }
    return false
}

// Assume points have all been evaluated, true if they have equal InvDivergence
func Uniform(context *NativeRegionRenderContext) bool {
    // If inverse divergence on all points is the same, no need to subdivide
    points := context.MandelbrotPoints()
    first := points[0].InverseDivergence()
    for _, p := range points[1:] {
        if p.InverseDivergence() != first {
            return false
        }
    }
    return true
}