package sharedregion

import (
	"fmt"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
)

type RenderCommand uint

const (
	ThreadRender = RenderCommand(iota)
	ThreadStop
)

type RenderInput struct {
	Command RenderCommand
	Regions []SharedRegionNumerics
}

type RenderOutput struct {
	UniformRegions []SharedRegionNumerics
	Children       []SharedRegionNumerics
	Members        []base.PixelMember
}

type RenderThread struct {
	ThreadId   uint
	InputChan  <-chan RenderInput
	OutputChan chan<- RenderOutput
	SharedConfig     SharedRegionConfig
	RegionConfig 	 region.RegionConfig
	BaseConfig	base.BaseConfig
	buffSize uint
	memberBuffSize uint
}

type RenderThreadFactory struct {
	count uint
	regionConfig region.RegionConfig
	sharedConfig SharedRegionConfig
}
// Easy method to create Render threads
func NewRenderThreadFactory(app RenderApplication) *RenderThreadFactory {
	return &RenderThreadFactory{
		count: 0,
		regionConfig: app.RegionConfig(),
		sharedConfig: app.SharedRegionConfig(),
	}
}

func (factory *RenderThreadFactory) Build(inputChan <-chan RenderInput, outputChan chan<- RenderOutput) RenderThread {
	collapseBound := factory.regionConfig.CollapseSize
	buffSize := factory.sharedConfig.BufferSize
	memberSize := collapseBound * collapseBound * buffSize
	thread := RenderThread{
		ThreadId:   factory.count,
		InputChan:  inputChan,
		OutputChan: outputChan,
		SharedConfig:     factory.sharedConfig,
		RegionConfig:	factory.regionConfig,
		buffSize: buffSize,
		memberBuffSize: memberSize,
	}
	factory.count++
	return thread
}

// Implements a single Render thread
func (thread *RenderThread) Run() {
	for {
		input := <-thread.InputChan
		switch input.Command {
		case ThreadRender:
			thread.OutputChan <- thread.Pass(input.Regions)
		case ThreadStop:
			return
		default:
			panic(fmt.Sprintf("Unknown Render Command in thread %v: %v",
				thread.ThreadId, input.Command))
		}
	}
}

// A pass through the region rendering process, comprising many steps
func (thread *RenderThread) Pass(Regions []SharedRegionNumerics) RenderOutput {
	output := thread.createRenderOutput()
	for _, region := range Regions {
		thread.Step(region, &output)
	}
	return output
}

// A single Render step
func (thread *RenderThread) Step(shared SharedRegionNumerics, output *RenderOutput) {
	// We are in a thread, so we must be sure that we have our own copy of the context
	shared.GrabThreadPrototype(thread.ThreadId)
	// We use proxies to share objects, so we've got to ensure we're using the correct local data
	shared.ClaimExtrinsics()

	baseConfig := thread.BaseConfig
	regionConfig := thread.RegionConfig
	iterateLimit := baseConfig.IterateLimit
	glitchSamples := regionConfig.GlitchSamples
	collapseBound := int(regionConfig.CollapseSize)

	if region.Collapse(shared, collapseBound) {
		points := SharedSequenceCollapse(shared, thread.ThreadId, iterateLimit)
		output.Members = append(output.Members, points...)
		return
	}

	if region.Subdivide(shared, iterateLimit, glitchSamples) {
		output.Children = append(output.Children, shared.SharedChildren()...)
		return
	}

	output.UniformRegions = append(output.UniformRegions, shared)
}

// Create a new output packet
func (thread *RenderThread) createRenderOutput() RenderOutput {
	return RenderOutput{
		// Q: How much memory should we allocate in advance?
		UniformRegions: make([]SharedRegionNumerics, 0, thread.buffSize),
		// What is the correct size for this?
		Members:  make([]base.PixelMember, 0, thread.memberBuffSize),
		Children: make([]SharedRegionNumerics, 0, thread.buffSize),
	}
}
