package sharedregion

import (
	"image"
	"math/rand"
	"testing"
	"time"
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

func TestNewRenderTracker(t *testing.T) {
	const jobCount = 5
	mock := &MockRenderApplication{}
	mock.SharedConfig.Jobs = uint32(jobCount)
	mock.SharedFactory = &MockFactory{}
	tracker := NewRenderTracker(mock)

	if !(mock.TSharedRegionConfig && mock.TDrawingContext) {
		t.Error("Expected methods not called on mock", mock)
	}

	if tracker == nil {
		t.Error("Expected tracker to be non-nil")
	}

	inputCount := len(tracker.input)
	if inputCount != jobCount {
		t.Error("Thread input channels of unexpected length: ", inputCount)
	}
	outputCount := len(tracker.output)
	if outputCount != jobCount {
		t.Error("Thread output channels of unexpected length", outputCount)
	}
	processCount := len(tracker.processing)
	if processCount != jobCount {
		t.Error("Thread processing tracker of unexpected length", processCount)
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
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	const jobCount = 3
	tracker := RenderTracker{
		jobs: jobCount,
		processing: make([]uint32, jobCount),
		input:      make([]chan RenderInput, jobCount),
	}

	// Create input channels
	for i := 0; i < jobCount; i++ {
		tracker.input[i] = make(chan RenderInput)
	}

	// Not zero-value input struct
	expect := RenderInput{Command: ThreadStop}

	// Pump input channels
	go func() {
		for range tracker.input {
			timeout(t, func() <-chan bool { return tracker.sendInput(expect) })
		}
	}()

	// Check each thread has input
	for _, threadInput := range tracker.input {
		timeout(t, func() <-chan bool {
			done := make(chan bool, 1)
			<-threadInput
			done <- true
			return done
		})
	}
}


// todo find a library that does this already
func timeout(t *testing.T, f func() <-chan bool) {
	timer := make(chan bool, 1)
	done := f()
	go func() {
		time.Sleep(1 * time.Second)
		timer <- true
	}()

	select {
	case <-done:
		return
	case <-timer:
		t.Error("Timed out")
	}
}

func TestTrackerRenderRegions(t *testing.T) {
	const jobCount = 1
	const threadBuff = 2
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
	startInputChan := make(chan RenderInput)
	startRegions := []SharedRegionNumerics{&MockNumerics{}}
	startTracker := RenderTracker{
		jobs: jobCount,
		buffer:     []SharedRegionNumerics{},
		processing: []uint32{0},
		input:      []chan RenderInput{startInputChan},
	}
	startTracker.renderRegions(startRegions)
	startOut := <-startTracker.input[0]
	if startOut.Command != ThreadRender && len(startOut.Regions) != 1 {
		t.Error("Read unexpected input:", startOut)
	}
	if len(startTracker.buffer) != 0 {
		t.Error("Tracker had unexpected buffer length:", startTracker)
	}

	// When busy, wait till the buffer fills
	busyInputChan := make(chan RenderInput)
	busyRegions := []SharedRegionNumerics{&MockNumerics{}}
	busyTracker := RenderTracker{
		jobs: jobCount,
		buffer:     []SharedRegionNumerics{},
		processing: []uint32{1},
		input:      []chan RenderInput{busyInputChan},
		config:     SharedRegionConfig{BufferSize: threadBuff},
	}
	busyTracker.renderRegions(busyRegions)
	if len(busyTracker.buffer) != 1 {
		t.Error("Tracker had unexpected buffer length")
	}
	busyTracker.renderRegions(busyRegions)
	out := <-busyTracker.input[0]
	if out.Command != ThreadRender && len(out.Regions) != 2 {
		t.Error("Read unexpected input:", out)
	}
	if len(busyTracker.buffer) != 0 {
		t.Error("Tracker had unexpected buffer length")
	}
}

func TestTrackerStep(t *testing.T) {
	const jobCount = 1
	child := &MockNumerics{}
	uniform := &MockNumerics{}
	member := base.PixelMember{I: 1, J: 2}
	out := RenderOutput{
		Children:       []SharedRegionNumerics{child},
		UniformRegions: []SharedRegionNumerics{uniform},
		Members:        []base.PixelMember{member},
	}

	// The tracker is not busy
	inputChan := make(chan RenderInput)
	tracker := RenderTracker{
		jobs: jobCount,
		buffer:     []SharedRegionNumerics{},
		processing: []uint32{0},
		input:      []chan RenderInput{inputChan},
		uniform:    []SharedRegionNumerics{},
		points:     []base.PixelMember{},
	}

	tracker.step(out)

	actual := <-tracker.input[0]
	if len(actual.Regions) != 1 {
		t.Error("Expected 1 region in the input channel but received:", actual)
	}

	if len(tracker.uniform) != 1 {
		t.Error("Expected 1 uniform region but tracker was:", tracker)
	}

	if len(tracker.points) != 1 {
		t.Error("Expected 1 base.PixelMember but tracker was:", tracker)
	}
}

func TestTrackerDraw(t *testing.T) {
	const iterateLimit = 255
	uniform := uniformer()
	point := base.PixelMember{I: 1, J: 2, Member: base.BaseMandelbrot{}}
	context := &draw.MockDrawingContext{
		Pic: image.NewNRGBA(image.ZR),
		Col: draw.NewRedscalePalette(iterateLimit),
	}
	tracker := RenderTracker{
		uniform: []SharedRegionNumerics{uniform},
		points:  []base.PixelMember{point},
		context:    context,
	}

	tracker.draw()

	if !(uniform.TRect && uniform.TRegionMember) {
		t.Error("Expected method not called on uniform region:", *uniform)
	}

	if !(context.TPicture && context.TColors) {
		t.Error("Expected method not called on drawing context")
	}
}


func sameInput(a RenderInput, b RenderInput) bool {
	return a.Command != b.Command && len(a.Regions) == len(b.Regions)
}