package sharedregion


type SharedRegionRenderStrategy RenderTracker

func Make(app RenderApplication) *SharedRegionRenderStrategy {
	tracker := NewRenderTracker(app)
    return (*SharedRegionRenderStrategy)(tracker)
}
