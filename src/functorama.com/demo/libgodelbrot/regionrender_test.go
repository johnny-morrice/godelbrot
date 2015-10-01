package libgodelbrot

import (
	"testing"
)

func BenchmarkRegionRender(b *testing.B) {
	config := DefaultConfig()
	redscale := NewRedscalePalette(DefaultIterations)
	pic := config.BlankImage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegionRenderImage(CreateContext(config, redscale, pic))
	}
}
