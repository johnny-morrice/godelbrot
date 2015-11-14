package sharedregion

import (
	"testing"
	"sync"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
)

const children = 4
const collapseCount = 1
const collapseSize = 10

type outputExpect struct {
    uniform  int
    children int
    members  int
}

func TestWorkerFactory(t *testing.T) {
	mock := &MockRenderApplication{}
	factory := NewWorkerFactory(mock, RenderOutput{})

	workers := []*Worker{factory.Build(), factory.Build()}

	if !mock.TSharedRegionConfig {
		t.Error("Mock did not receive expected method call")
	}

	for i, th := range workers {
		if th.WorkerId != uint16(i) {
			t.Error("Worker", i, "had incorrect WorkerId: ", th)
		}
	}
}

func TestWorkerRun(t *testing.T) {
	worker := newWorker()
	const uniformLength = 1

	go func() {
		worker.InputChan <- RenderInput{
			Region: uniformer(),
		}
		close(worker.InputChan)
	}()

	go func() {
		for range worker.WaitingChan {
			// pump the channel
		}
		worker.Output.Close()
	}()

	go worker.Run()

	workerOutputCheck(t, worker, outputExpect{1, 0, 0})
}

func TestWorkerStep(t *testing.T) {
	coll := collapser()
	subd := subdivider()
	uni := uniformer()
	collapsed := workerStep(coll)
	subdivided := workerStep(subd)
	uniformed := workerStep(uni)

	for i, mock := range []*MockNumerics{coll, subd, uni} {
		if !stepOkayGeneral(mock) {
			t.Error("General case methods not called on region", i, mock)
		}
	}

	if !(coll.TSharedRegionSequence) {
		t.Error("Expected methods not called on collapse region:", coll)
	}

	if !(subd.TSplit && subd.TSharedChildren && subd.TEvaluateAllPoints) {
		t.Error("Expected methods not called on subdivided region:", subd)
	}

	if !(uni.TOnGlitchCurve && uni.TEvaluateAllPoints)  {
		t.Error("Expected methods not called on uniform region:", uni)
	}

	workers := []*Worker{collapsed, subdivided, uniformed}

	go func() {
		for _, worker := range workers {
			worker.Hold.Wait()
			worker.Output.Close()
		}
	}()

	workerOutputCheck(t, collapsed, outputExpect{0, 0, collapseCount})
	workerOutputCheck(t, subdivided, outputExpect{0, children, 0})
	workerOutputCheck(t, uniformed, outputExpect{1, 0, 0})
}

func stepOkayGeneral(mock *MockNumerics) bool {
	okay := mock.TRect && mock.TGrabWorkerPrototype
	okay = okay && mock.TClaimExtrinsics
	return okay
}

func collapser() *MockNumerics {
	mock := mocker(region.CollapsePath)
	return mock
}

func subdivider() *MockNumerics {
	mock := mocker(region.SubdividePath)
	children := make([]*MockNumerics, children)

	for i := 0; i < len(children); i++ {
		children[i] = uniformer()
	}

	mock.SharedMockChildren = children
	return mock
}

func uniformer() *MockNumerics {
	return mocker(region.UniformPath)
}

func mocker(path region.RegionType) *MockNumerics {
	mock := &MockNumerics{}
	mock.Path = path
	mockSequence := &MockSequence{}
	mock.SharedMockSequence = mockSequence
	mock.AppCollapseSize = collapseSize
	return mock
}

func workerOutputCheck(t *testing.T, worker *Worker, expect outputExpect) {
	actualUniformCount := 0
	actualMemberCount := 0
	actualChildCount := 0

	hold := sync.WaitGroup{}
	const outputFieldCount = 3
	hold.Add(outputFieldCount)

	// Drain output channels
	go func() {
		for range worker.Output.UniformRegions {
			actualUniformCount++
		}
		hold.Done()
	}()

	go func() {
		for range worker.Output.Members {
			actualMemberCount++
		}
		hold.Done()
	}()

	go func() {
		for range worker.Output.Children {
			actualChildCount++
		}
		hold.Done()
	}()

	hold.Wait()

	okay := actualMemberCount == expect.members
	okay = okay && actualChildCount == expect.children
	okay = okay && actualUniformCount == expect.uniform
	if !okay {
		t.Error("Expected output counts", expect, "but received (",
			actualUniformCount, actualChildCount, actualMemberCount, ")")
	}
}

func workerStep(numerics SharedRegionNumerics) *Worker {
	worker := newWorker()
	worker.Step(numerics)
	return worker
}

func newWorker() *Worker {
	return &Worker{
		InputChan:  make(chan RenderInput),
		WaitingChan: make(chan bool),
		Output: RenderOutput{
			UniformRegions: make(chan SharedRegionNumerics),
			Children: make(chan SharedRegionNumerics),
			Members: make(chan base.PixelMember),
		},
		RegionConfig: region.RegionConfig{
			CollapseSize: collapseSize,
		},
	}
}
