package libgodelbrot

import (
    "testing"
)

func TestNewRenderTracker(t *testing.T) {
    jobCount := 5
    mock := mockRenderApplication{}
    mock.concurrentConfig.Jobs = jobCount
    tracker := NewRenderTracker(mock)

    if !(mock.tConcurrentConfig && mock.tDrawingContext) {
        t.Error("Expected methods not called on mock", mock)
    }

    if tracker == nil {
        t.Error("Expected tracker to be non-nil")
    }

    threadData := []interface{}{
        tracker.input,
        tracker.output,
        tracker.processing
    }

    for i, threadSlice := range threadData {
        actualCount := len(threadSlice)
        if actualCount != jobCount {
            t.Error("Data item", i, "had unexpected length: ", actualCount)
        }
    }
}