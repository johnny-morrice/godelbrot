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

func (output RenderOutput) Close() {
	close(output.UniformRegions)
	close(output.Children)
	close(output.Members)
}

type Worker struct {
	WorkerId   uint16
	InputChan  chan RenderInput
	Output RenderOutput
	WaitingChan chan bool
	Close chan bool
	SharedConfig     SharedRegionConfig
	RegionConfig 	 region.RegionConfig
	Hold sync.WaitGroup
}

type WorkerFactory struct {
	count uint16
	regionConfig region.RegionConfig
	sharedConfig SharedRegionConfig
	baseConfig base.BaseConfig
	output RenderOutput
}

// Easy method to create Render workers
func NewWorkerFactory(app RenderApplication, outputChannels RenderOutput) *WorkerFactory {
	return &WorkerFactory{
		count: 0,
		regionConfig: app.RegionConfig(),
		sharedConfig: app.SharedRegionConfig(),
		output: outputChannels,
	}
}

func (factory *WorkerFactory) Build() *Worker {
	worker := &Worker{
		WorkerId:   factory.count,
		InputChan:  make(chan RenderInput),
		WaitingChan: make(chan bool),
		Close: make(chan bool),
		SharedConfig:     factory.sharedConfig,
		RegionConfig:	factory.regionConfig,
		Output: factory.output,
	}
	factory.count++
	return worker
}


// TODO introduce dead-wait status to indicate that data is flushed

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
		case input := <-worker.InputChan:
			// Enter busy state
			busy = true
			worker.WaitingChan<- false
			worker.Step(input.Region)
		case <-worker.Close:
			worker.closeChannels()
			return
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

	regionConfig := worker.RegionConfig
	collapseBound := int(regionConfig.CollapseSize)

	if region.Collapse(shared, collapseBound) {
		points := SharedSequenceCollapse(shared, worker.WorkerId)
		worker.Hold.Add(1)
		go func() {
			for _, p := range points {
				worker.Output.Members<- p
			}
			worker.Hold.Done()
		}()
		return
	}

	if region.Subdivide(shared) {
		children := shared.SharedChildren()
		worker.Hold.Add(len(children))
		for _, child := range children {
			go func(spawn SharedRegionNumerics) {
				worker.Output.Children<- spawn
				worker.Hold.Done()
			}(child)
		}
		return
	}

	worker.Hold.Add(1)
	go func() {
		worker.Output.UniformRegions<- shared
		worker.Hold.Done()
	}()
}

func (worker *Worker) closeChannels() {
	worker.Hold.Wait()
	close(worker.WaitingChan)
	close(worker.InputChan)
}