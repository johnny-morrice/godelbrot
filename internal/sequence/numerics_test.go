package sequence

import (
    "testing"
    "github.com/johnny-morrice/godelbrot/internal/draw"
)

func TestImageSequence(t *testing.T) {
    const iterateLimit = 10
    context := draw.NewMockDrawingContext(iterateLimit)
    numerics := &MockNumerics{}
    ImageSequence(numerics, context)

    if false {
        t.Error("Since the mock context does not have a real image attached",
            "all this test really does is prove that the mechanism works without crashing out")
    }
}