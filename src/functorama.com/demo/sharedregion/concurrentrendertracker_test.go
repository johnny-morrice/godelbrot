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

func TestStashUniforms(t *testing.T) {
	mockA := &MockNumerics{}
	mockB := &MockNumerics{}
	uniforms := []SharedRegionNumerics{mockA, mockB}

	tracker := &RenderTracker{
		uniformChan: make(chan SharedRegionNumerics),
	}

	go tracker.stashUniforms(uniforms)

	actualAGen := <-tracker.uniformChan
	actualBGen := <-tracker.uniformChan

	actualA := actualAGen.(*MockNumerics)
	actualB := actualBGen.(*MockNumerics)

	if actualA != mockA {
		t.Error("Expected", mockA, "but received", actualA)
	}

	if actualB != mockB {
		t.Error("Expected", mockB, "but received", actualB)
	}
}

func TestStashMembers(t *testing.T) {
	pointA := base.PixelMember{I: 1}
	pointB := base.PixelMember{I: 2}
	points := []base.PixelMember{pointA, pointB}

	tracker := &RenderTracker{
		memberChan: make(chan base.PixelMember),
	}

	go tracker.stashMembers(points)

	actualA := <-tracker.uniformChan
	actualB := <-tracker.uniformChan

	if actualA != pointA {
		t.Error("Expected", pointA, "but received", actualA)
	}

	if actualB != pointB {
		t.Error("Expected", pointB, "but received", actualB)
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

	tracker.draw()

	if !(uniform.TRect && uniform.TRegionMember) {
		t.Error("Expected method not called on uniform region:", *uniform)
	}

	if !(context.TPicture && context.TColors) {
		t.Error("Expected method not called on drawing context")
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

func sameInput(a RenderInput, b RenderInput) bool {
	return a.Command != b.Command && len(a.Regions) == len(b.Regions)
}