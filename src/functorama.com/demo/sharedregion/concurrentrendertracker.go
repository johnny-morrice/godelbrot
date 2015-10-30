package sharedregion

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
	"functorama.com/demo/region"
)

type RenderTracker struct {
	// nth element incremented when nth thread processing
	processing []uint32
	// input channels to threads
	input []chan RenderInput
	// output channels to threads
	output []chan RenderOutput
	// round robin input scheduler
	nextThread int
	// Local buffer of unprocessed regions
	buffer []SharedRegionNumerics
	// Concurrent render config
	config SharedRegionConfig
	// drawing context for drawing onto image
	context draw.DrawingContext
	// Processed uniform regions
	uniform []SharedRegionNumerics
	// Processed Mandelbrot points
	points []base.PixelMember
	// Thread factory
	factory *RenderThreadFactory
	// Initial region
	initialRegion SharedRegionNumerics
}

func NewRenderTracker(app RenderApplication) *RenderTracker {
	config := app.SharedRegionConfig()
	tracker := RenderTracker{
		processing: make([]uint32, config.Jobs),
		input:      make([]chan RenderInput, config.Jobs),
		output:     make([]chan RenderOutput, config.Jobs),
		buffer:     newBuffer(config.BufferSize),
		nextThread: 0,
		config:     config,
		context:       app.DrawingContext(),
		uniform:    make([]SharedRegionNumerics, base.AllocMedium),
		points:     make([]base.PixelMember, base.AllocLarge),
		factory:	NewRenderThreadFactory(app),
		initialRegion: app.SharedRegionFactory().Build(),
	}
	for i := 0; i < int(config.Jobs); i++ {
		tracker.processing[i] = 0
		inputChan := make(chan RenderInput, base.AllocSmall)
		outputChan := make(chan RenderOutput, base.AllocSmall)
		tracker.input[i] = inputChan
		tracker.output[i] = outputChan
	}
	return &tracker
}

// True if at least one thread is processing
func (tracker *RenderTracker) busy() bool {
	for _, wait := range tracker.processing {
		if wait > 0 {
			return true
		}
	}
	return false
}

// Render a set of regions using the next thread in the round robin scheme
func (tracker *RenderTracker) renderRegions(regions []SharedRegionNumerics) {

	tracker.buffer = append(tracker.buffer, regions...)

	if len(tracker.buffer) == 0 {
		return
	}

	input := RenderInput{Command: ThreadRender}

	// If we're busy, wait for buffer to fill before sending
	if tracker.busy() {
		if len(tracker.buffer) >= int(tracker.config.BufferSize) {
			input.Regions = tracker.buffer
			tracker.cleanBuffer()
			tracker.sendInput(input)
		}
	} else {
		input.Regions = tracker.buffer
		tracker.cleanBuffer()
		tracker.sendInput(input)
	}
}

// Send input and mark as busy
func (tracker *RenderTracker) sendInput(input RenderInput) {
	threadIndex := tracker.nextThread
	tracker.input[threadIndex] <- input
	tracker.processing[threadIndex]++

	// Rount robin
	tracker.nextThread++
	if tracker.nextThread >= len(tracker.processing) {
		tracker.nextThread = 0
	}
}

// A single step in render tracking
func (tracker *RenderTracker) step(output RenderOutput) {

	// Give more work to the threads
	tracker.renderRegions(output.Children)

	// Stash the completed areas
	tracker.uniform = append(tracker.uniform, output.UniformRegions...)
	tracker.points = append(tracker.points, output.Members...)
}

// draw to the image
func (tracker *RenderTracker) draw() {
	for _, uniform := range tracker.uniform {
		uniform.ClaimExtrinsics()
		region.DrawUniform(tracker.context, uniform)
	}

	// We do not need to claim any extrinsics here, because we are merely drawing a render result
	// that requires no extra context
	for _, member := range tracker.points {
		draw.DrawPoint(tracker.context, member)
	}
}

func (tracker *RenderTracker) cleanBuffer() {
	tracker.buffer = newBuffer(tracker.config.BufferSize)
}

func newBuffer(bufferSize uint) []SharedRegionNumerics {
	return make([]SharedRegionNumerics, 0, bufferSize)
}

// How do we test Render?

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
	// Launch threads
	for i := uint32(0); i < tracker.config.Jobs; i++ {
		thread := tracker.factory.Build(tracker.input[i], tracker.output[i])
		go thread.Run()
	}

	firstBatch := []SharedRegionNumerics{tracker.initialRegion}
	tracker.renderRegions(firstBatch)

	for tracker.busy() {
		for i, outputChan := range tracker.output {
			select {
			case output := <-outputChan:
				tracker.processing[i]--
				tracker.step(output)
			default:
				// Wait till next input
			}
		}
	}

	// Shut down threads
	for _, input := range tracker.input {
		input <- RenderInput{Command: ThreadStop}
	}

	tracker.draw()
}
