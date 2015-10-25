package libgodelbrot

type RegionNumericsFactory interface{
    Region() RegionNumerics
}

// RegionNumerics provides rendering calculations for the "region" render strategy.
type RegionNumerics interface {
    OpaqueProxyFlyweight
    Rect() image.Rectangle
    EvaluateAllPoints(iterateLimit int)
    Split()
    OnGlitchCurve(iterateLimit uint8, glitchSamples uint) bool
    MandelbrotPoints() []MandelbrotMember
    RegionMember()
    Subdivide() bool
    Children() []RegionNumerics
    RegionalSequenceNumerics() SequentialNumerics
}


// RenderSequentialRegion takes a RegionNumerics but renders the region in a sequential
// (column-wise) manner
func RenderSequentialRegion(numerics RegionNumerics) {
    smallNumerics := numerics.RegionalSequenceNumerics()
    smallNumerics.ClaimExtrinsics()
    SequentialRenderImage(smallNumerics)
}

// SequenceCollapse is analogous to RenderSequentialRegion, but it returns the Mandelbrot render
// results rather than drawing them to the image.
func SequenceCollapse(numerics RegionNumerics) []PixelMember {
    sequence := region.RegionalSequenceNumerics()
    sequence.ClaimExtrinsics()
    sequence.MemberCaptureSequencer()
    MandelbrotSequence(smallConfig)
    return sequence.CapturedMembers()
}

// Subdivide takes a RegionNumerics and tries to split the region into subregions.  It returns true
// if the subdivision occurred.  The subdivision won't occur if the region is Uniform or in an area
// where glitches are likely.
func Subdivide(numerics RegionNumerics) bool {
    numerics.EvaluateAllPoints()

    if !Uniform(numerics) || numerics.OnGlitchCurve() {
        numerics.Split()
        return true
    }
    return false
}

// Uniform returns true if the region has the same Mandelbrot escape value across its bounds.
// Note this function performs no evaluation of membership.
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

// Collapse returns true if the region is below the necessary size for subdivision
func Collapse(split RegionNumerics, collapseSize int) bool {
    rect := split.Rect()
    return rect.Dx() <= collapseSize || rect.Dy() <= collapseSize
}