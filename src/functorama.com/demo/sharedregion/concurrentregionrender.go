package sharedregion

import (
    "image"
)

type SharedRegionRenderStrategy RenderTracker

func Make(app RenderApplication) *SharedRegionRenderStrategy {
    return (*SharedRegionRenderStrategy)(NewRenderTracker(app))
}

func (srrs *SharedRegionRenderStrategy) Render() (*image.NRGBA, error) {
    track := (*RenderTracker)(srrs)
    track.Render()
    return track.context.Picture(), nil
}