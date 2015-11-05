package sharedregion


type SharedRegionRenderStrategy RenderTracker

func NewSharedRegionRenderer(app RenderApplication) *SharedRegionRenderStrategy {
	tracker := NewRenderTracker(app)
    return (*SharedRegionRenderStrategy)(tracker)
}
