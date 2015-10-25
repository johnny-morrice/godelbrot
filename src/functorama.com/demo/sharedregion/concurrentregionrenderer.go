package sharedregion


func NewSharedRegionRenderer(app RenderApplication) *RenderTracker {
	return NewRenderTracker(app)
}
