package libgodelbrot

import (
	"testing"
)

func BenchmarkSequentialRender(b *testing.B) {
	config := DefaultConfig()
	redscale := NewRedscalePalette(DefaultIterations)
	pic := config.BlankImage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SequentialRenderImage(CreateContext(config, redscale, pic))
	}
}
