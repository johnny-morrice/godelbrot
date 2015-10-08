package libgodelbrot

// These datatypes are shared between the concurrent RenderTracker and 
// RegionRenderThreads

type renderCommand uint

const (
    render = renderCommand(iota)
    stop   = renderCommand(iota)
)

type renderInput struct {
    command renderCommand
    regions []RegionRenderContext
}

type pixelMember struct {
    i int
    j int
    MandelbrotMember
}

type renderOutput struct {
    uniformRegions []RegionRenderContext
    children       []RegionRenderContext
    members        []pixelMember
}