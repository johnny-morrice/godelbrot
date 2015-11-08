package sharedregion

import (
	"sync"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
)

type RenderInput struct {
	Region SharedRegionNumerics
}

type RenderOutput struct {
	UniformRegions chan SharedRegionNumerics
	Children       chan SharedRegionNumerics
	Members        chan base.PixelMember
}

type Worker struct {
	WorkerId   uint16
	InputChan  chan RenderInput
	Output RenderOutput
	WaitingChan chan bool
	SharedConfig     SharedRegionConfig
	RegionConfig 	 region.RegionConfig
	BaseConfig	base.BaseConfig
	hold sync.WaitGroup
}

type WorkerFactory struct {
	count uint16
	regionConfig region.RegionConfig
	sharedConfig SharedRegionConfig
}

// Easy method to create Render workers
func NewWorkerFactory(app RenderApplication) *WorkerFactory {
	return &WorkerFactory{
		count: 0,
		regionConfig: app.RegionConfig(),
		sharedConfig: app.SharedRegionConfig(),
	}
}

func (factory *WorkerFactory) Build(outputChannels RenderOutput) *Worker {
	worker := &Worker{
		WorkerId:   factory.count,
		InputChan:  make(chan RenderInput),
		WaitingChan: make(chan bool),
		SharedConfig:     factory.sharedConfig,
		RegionConfig:	factory.regionConfig,
		Output: outputChannels,
	}
	factory.count++
	return worker
}

// Implements a single Render worker
func (worker *Worker) Run() {
	busy := true
	for {
		// Enter wait state
		if busy {
			worker.WaitingChan<- true
			busy = false
		}
		select {
		case input, ok := <-worker.InputChan:
			if ok {
				// Enter busy state
				busy = true
				worker.WaitingChan<- false
				worker.Step(input.Region)
			} else {
				worker.closeChannels()
				return
			}
		default:
			continue
		}

	}
}

// A single Render step
func (worker *Worker) Step(shared SharedRegionNumerics) {
	// We are in a worker, so we must be sure that we have our own copy of the context
	shared.GrabWorkerPrototype(worker.WorkerId)
	// We use proxies to share objects, so we've got to ensure we're using the correct local data
	shared.ClaimExtrinsics()

	baseConfig := worker.BaseConfig
	regionConfig := worker.RegionConfig
	iterateLimit := baseConfig.IterateLimit
	glitchSamples := regionConfig.GlitchSamples
	collapseBound := int(regionConfig.CollapseSize)

	if region.Collapse(shared, collapseBound) {
		points := SharedSequenceCollapse(shared, worker.WorkerId, iterateLimit)
		worker.hold.Add(len(points))
		for _, point := range points {
			go func() {
				worker.Output.Members<- point
				worker.hold.Done()
			}()
		}
		return
	}

	if region.Subdivide(shared, iterateLimit, glitchSamples) {
		children := shared.SharedChildren()
		worker.hold.Add(len(children))
		for _, child := range children {
			go func() {
				worker.Output.Children<- child
				worker.hold.Done()
			}()
		}
		return
	}

	worker.hold.Add(1)
	go func() {
		worker.Output.UniformRegions<- shared
		worker.hold.Done()
	}()
}

func (worker *Worker) closeChannels() {
	worker.hold.Wait()
	close(worker.Output.Children)
	close(worker.Output.UniformRegions)
	close(worker.Output.Members)
	close(worker.WaitingChan)
}