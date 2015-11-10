package sharedregion

import (
	"image"
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

	workerCount := len(tracker.workers)
	if workerCount != jobCount {
		t.Error("Expected", jobCount, "workers but received", workerCount)
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
		uniformChan: make(chan SharedRegionNumerics),
		memberChan: make(chan base.PixelMember),
		context:    context,
	}

	go func() {
		tracker.uniformChan<- uniform
		close(tracker.uniformChan)
	}()
	go func() {
		tracker.memberChan<- point
		close(tracker.memberChan)
	}()

	drawPackets := tracker.syncDrawing
	tracker.draw(drawPackets)

	if !(uniform.TRect && uniform.TRegionMember) {
		t.Error("Expected method not called on uniform region:", *uniform)
	}

	if !(context.TPicture && context.TColors) {
		t.Error("Expected method not called on drawing context")
	}
}

func TestTrackerCirculate(t *testing.T) {
	tracker := &RenderTracker{
		workersDone: make(chan bool),
		childChan: make(chan SharedRegionNumerics),
		schedule: make(chan chan RenderInput),
	}

	expectedRegion := &MockNumerics{}

	// Feed input
	workerInput := make(chan RenderInput)
	go func() {
		tracker.schedule<- workerInput	
	}()
	go func() {
		tracker.childChan<- expectedRegion
	}()
	done := make(chan bool, 1)
	go func() {
		tracker.circulate()
		done<- true
	}

	// Test input
	abstractNumerics <-workerInput
	actualRegion := abstractNumerics.(*MockNumerics)
	if actualRegion != expectedRegion {
		t.Error("Expected", expectedRegion,
			"but received", actualRegion)
	}

	// Test shutdown
	go func() {
		tracker.workersDone<- true
		timeout(t, func() { return done })
	}
}

func TestTrackerScheduleWorkers(t *testing.T) {
	const jobCount = 2
	tracker := &RenderTracker{
		workers: make([]Worker, jobCount),
		schedule: make(chan chan RenderInput),
		stateChan: make(chan workerState),
		workerOutput: RenderOutput{
			UniformRegions: make(chan SharedRegionNumerics),
			Children: make(chan SharedRegionNumerics),
			Members: make(chan base.PixelMember),
		},
	}

	app := &MockRenderApplication{}
	factory := NewWorkerFactory(app)

	workerA := factory.Build(tracker.workerOutput)
	workerB := factory.Build(tracker.workerOutput)

	tracker.workers = []Worker{workerA, workerB}

	// Run schedule process
	stop := tracker.scheduleWorkers()

	// Test input scheduling
	go func() {
		workerA.ReadyChan<- true
	}
	actualA := <-tracker.schedule

	if actualA != workerA.inputChan {
		t.Error("Expected", workerA.inputChan,
			"but received", actualA)
	}

	go func() {
		workerB.ReadyChan<- true
	}
	actualB := <-tracker.schedule

	if actualB != workerB.inputChan {
		t.Error("Expected", workerB.inputChan,
			"but received", actualB)
	}

	stop<- true
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

func sameInput(a RenderInput, b RenderInput) bool {
	return a.Command != b.Command && len(a.Regions) == len(b.Regions)
}