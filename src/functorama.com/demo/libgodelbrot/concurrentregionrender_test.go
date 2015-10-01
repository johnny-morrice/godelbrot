package libgodelbrot 

import (
    "testing"
)

func BenchmarkConcurrentRegionRender(b *testing.B) {
    config := DefaultConfig()
    redscale := NewRedscalePalette(DefaultIterations)
    pic := config.BlankImage()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ConcurrentRegionRenderImage(CreateContext(config, redscale, pic))
    }
}