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
	Hold sync.WaitGroup
	sizelim int
}

type WorkerFactory struct {
	count uint16
	regionConfig region.RegionConfig
	baseConfig base.BaseConfig
	output RenderOutput
}

// Easy method to create Render workers
func NewWorkerFactory(app RenderApplication, outputChannels RenderOutput) *WorkerFactory {
	return &WorkerFactory{
		count: 0,
		regionConfig: app.RegionConfig(),
		output: outputChannels,
	}
}

func (factory *WorkerFactory) Build() *Worker {
	worker := &Worker{
		WorkerId:   factory.count,
		InputChan:  make(chan SharedRegionNumerics),
		sizelim:	int(factory.regionConfig.CollapseSize),
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
func (worker *Worker) Step(reg SharedRegionNumerics) {
	reg.GrabWorkerPrototype(worker.WorkerId)
	reg.ClaimExtrinsics()

	if region.Collapse(reg, worker.sizelim) {
		px := SharedSequenceCollapse(reg, worker.WorkerId)
		worker.Hold.Add(1)
		go func() {
			out := WorkerPixelOut{Id: worker.WorkerId, Points: px}
			worker.Output.Members<- out
			worker.Hold.Done()
		}()
		return
	}

	if region.Subdivide(reg) {
		children := reg.SharedChildren()
		worker.Hold.Add(1)
		go func() {
			out := WorkerChildrenOut{Id: worker.WorkerId, Children: children}
			worker.Output.Children<- out
			worker.Hold.Done()
		}()
		return
	}

	worker.Hold.Add(1)
	go func() {
		out := WorkerRegionOut{Id: worker.WorkerId, Region: reg}
		worker.Output.UniformRegions<- out
		worker.Hold.Done()
	}()
}