package sharedregion

import (
	"sync"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
)

type RenderOutput struct {
	UniformRegions chan WorkerRegionOut
	Children       chan WorkerChildrenOut
	Members        chan WorkerPixelOut
}

func (output RenderOutput) Close() {
	close(output.UniformRegions)
	close(output.Children)
	close(output.Members)
}

type WorkerChildrenOut struct {
	Id uint16
	Children []SharedRegionNumerics
}

type WorkerRegionOut struct {
	Id uint16
	Region SharedRegionNumerics
}

type WorkerPixelOut struct {
	Id uint16
	Points []base.PixelMember
}

type Worker struct {
	WorkerId   uint16
	InputChan  chan SharedRegionNumerics
	Output RenderOutput
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
		InputChan:  make(chan SharedRegionNumerics),
		SharedConfig:     factory.sharedConfig,
		RegionConfig:	factory.regionConfig,
		Output: factory.output,
	}
	factory.count++
	return worker
}


func (worker *Worker) Run() {
	for i := range worker.InputChan {
		worker.Step(i)
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
			pixels := WorkerPixelOut{Id: worker.WorkerId, Points: points}
			worker.Output.Members<- pixels
			worker.Hold.Done()
		}()
		return
	}

	if region.Subdivide(shared) {
		children := shared.SharedChildren()
		worker.Hold.Add(1)
		go func() {
			reg := WorkerChildrenOut{Id: worker.WorkerId, Children: children}
			worker.Output.Children<- reg
			worker.Hold.Done()
		}()
		return
	}

	worker.Hold.Add(1)
	go func() {
		reg := WorkerRegionOut{Id: worker.WorkerId, Region: shared}
		worker.Output.UniformRegions<- reg
		worker.Hold.Done()
	}()
}