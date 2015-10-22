package libgodelbrot

import (
	"math/rand"
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
		tracker.processing,
	}

	for i, threadSlice := range threadData {
		actualCount := len(threadSlice)
		if actualCount != jobCount {
			t.Error("Data item", i, "had unexpected length: ", actualCount)
		}
	}
}

func TestTrackerBusy(t *testing.T) {
	jobCount := 5
	tracker := RenderTracker{
		processing: make([]uint32, jobCount),
	}

	if tracker.busy() {
		t.Error("New tracker should not be busy")
	}

	for i := 0; i < jobCount; i++ {
		for j := 0; j < jobCount; j++ {
			tracker.processing[j] = 0
		}
		tracker.processing[i] = rand.Uint32()
		if !tracker.busy() {
			t.Error("Expected tracker to be busy. Tracker processing: ", tracker.processing)
		}
	}
}

func TestTrackerSendInput(t *testing.T) {
	jobCount := 3
	tracker := RenderTracker{
		processing: make([]uint32, jobCount),
		input:      make([]chan renderInput, jobCount),
	}

	// Create input channels
	for i := 0; i < jobCount; i++ {
		tracker.input[i] = make(chan renderInput)
	}

	// Not zero-value input struct
	sampleInput := renderInput{command: stop}

	// Pump input channels
	for i := 0; i < jobCount; i++ {
		if tracker.nextThread != i {
			t.Error("Unexpected target thread")
		}
		tracker.sendInput(sampleInput)
		out := <-tracker.input[i]
		if out != sampleInput {
			t.Error("On channel", i, "expected input", sampleInput, "but received", out)
		}
	}

	// Ensure thread wrap around
	if tracker.nextThread != 0 {
		t.Error("Unexpected target thread after wrap")
	}
	tracker.sendInput(sampleInput)
	out := <-tracker.input[i]
	if out != sampleInput {
		t.Error("On channel 0 after wrap expected input", sampleInput, "but received", out)
	}
}

func TestTrackerRenderRegions(t *testing.T) {
	// This function does not reference particular threads
	const jobCount = 1

	// Check trivial path
	zeroRegions := []SharedRegionNumerics{}
	zeroTracker := RenderTracker{
		buffer: []SharedRegionNumerics{},
	}
	zeroTracker.renderRegions(zeroRegions)
	if len(zeroTracker.buffer) > 0 {
		t.Error("Unexpected buffer growth:", zeroTracker)
	}

	// When not busy, just send whatever we have to the next thread
	startRegions := []SharedRegionNumerics{mockSharedRegionNumerics{}}
	startTracker := RenderTracker{
		buffer:     []SharedRegionNumerics{},
		processing: []uint32{0},
		input:      []chan renderInput{make(chan renderInput)},
	}
	startTracker.renderRegions(startRegions)
	startOut := <-tracker.input[0]
	if startOut.command != render && len(startOut.regions) != 1 {
		t.Error("Read unexpected input:", startOut)
	}
	if len(startTracker.buffer) != 0 {
		t.Error("Tracker had unexpected buffer length:", startTracker)
	}

	// When busy, wait till the buffer fills
	busyRegions := []SharedRegionNumerics{mockSharedRegionNumerics{}}
	busyTracker := RenderTracker{
		buffer:     []SharedRegionNumerics{},
		processing: []uint32{1},
		input:      []chan renderInput{make(chan renderInput)},
		config:     ConcurrentRenderConfig{BufferSize: 2},
	}
	busyTracker.renderRegions(busyRegions)
	if len(busyTracker.buffer) != 1 {
		t.Error("Tracker had unexpected buffer length")
	}
	busyTracker.renderRegions(busyRegions)
	out := <-tracker.input[0]
	if out.command != render && len(out.regions) != 2 {
		t.Error("Read unexpected input:", out)
	}
	if len(busyTracker.buffer) != 0 {
		t.Error("Tracker had unexpected buffer length")
	}
}

func TestTrackerStep(t *testing.T) {
	
}
