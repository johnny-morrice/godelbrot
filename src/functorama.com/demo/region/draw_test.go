package region

import (
    "image"
    "testing"
    "functorama.com/demo/draw"
)

func TestDrawUniform(t *testing.T) {
    mockDraw := draw.MockDrawingContext{
        Pic: image.NewNRGBA(image.ZR),
        Col: draw.NewRedscalePalette(255),
    }
    mockRegion := mockRegionNumerics{}
    DrawUniform(mockDraw, mockRegion)

    if !(mockDraw.TPicture && mockDraw.TColors) {
        t.Error("Expected method not called on mock drawing context:", mockDraw)
    }

    if !(mockRegion.tRegionMember && mockRegion.tRect) {
        t.Error("Expected method not called on mock region numerics:", mockRegion)
    }
}