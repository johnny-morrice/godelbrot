package sharedregion

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
	"functorama.com/demo/region"
)

type RenderTracker struct {
	// Number of jobs
	jobs int
	working []chan bool
	ready []chan bool
	// schedule channel
	schedule chan chan<- RenderInput
	// input channels to threads
	input []chan RenderInput
	// output channels to threads
	output []chan RenderOutput
	// Concurrent render config
	config SharedRegionConfig
	// drawing context for drawing onto image
	context draw.DrawingContext
	uniformChan chan SharedRegionNumerics
	memberChan chan base.PixelMember
	// Thread factory
	factory *RenderThreadFactory
	// Initial region
	initialRegion SharedRegionNumerics
}

func NewRenderTracker(app RenderApplication) *RenderTracker {
	config := app.SharedRegionConfig()
	tracker := RenderTracker{
		jobs: int(config.Jobs),
		working: make([]chan bool, config.Jobs),
		ready: make([]chan bool, config.Jobs),
		input:      make([]chan RenderInput, config.Jobs),
		output:     make([]chan RenderOutput, config.Jobs),
		config:     config,
		context:       app.DrawingContext(),
		memberChan: make(chan base.PixelMember),
		uniformChan: make(chan SharedRegionNumerics),
		factory:	NewRenderThreadFactory(app),
		initialRegion: app.SharedRegionFactory().Build(),
	}

	for i := 0; i < int(config.Jobs); i++ {
		readyChan := make(chan bool, 1)
		workingChan := make(chan bool, 1)
		inputChan := make(chan RenderInput)
		outputChan := make(chan RenderOutput)
		tracker.working[i] = workingChan
		tracker.ready[i] = readyChan
		tracker.input[i] = inputChan
		tracker.output[i] = outputChan
	}
	return &tracker
}

// A single step in render tracking
func (tracker *RenderTracker) step(output RenderOutput) {

	// Give more work to the threads
	tracker.work(output.Children)

	// Stash the completed areas
	go func() {
		for _, uni := range output.UniformRegions {
			tracker.uniformChan<- uni
		}
	}()

	go func() {
		for _, member := range output.Members {
			tracker.memberChan<- member
		}
	}()
}

// draw to the image
func (tracker *RenderTracker) draw() {
	for uniform := range tracker.uniformChan {
		uniform.ClaimExtrinsics()
		region.DrawUniform(tracker.context, uniform)
	}

	// We do not need to claim any extrinsics here, because we are merely drawing a render result
	// that requires no extra context
	for member := range tracker.memberChan {
		draw.DrawPoint(tracker.context, member)
	}
}

func newBuffer(bufferSize uint) []SharedRegionNumerics {
	return make([]SharedRegionNumerics, 0, bufferSize)
}

func (tracker *RenderTracker) inputSchedule() (chan<- bool, <-chan bool) {
	stop := make(chan bool)
	busyChan := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				// Nothing here
			}
			busy := false
			for i, readyChan := range tracker.ready {
				select {
				case <-readyChan:
					tracker.schedule<- tracker.input[i]
				default:
					busy = true
					continue
				}
			}
			go func() { busyChan<- busy }()
		}
	}()

	return stop, busyChan
}

func (tracker *RenderTracker) circulate() (chan<- bool, <-chan bool) {
	stop := make(chan bool)
	waitingChan := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				// Nothing here
			}
			waiting := false
			for i, workingChan := range tracker.working {
				select {
				case <-workingChan:
					// Await result
					go func() {
						result := <- tracker.output[i]
						tracker.step(result)
					}()
					waiting = true
				default:
					continue
				}
			}
			go func() {waitingChan<- waiting}()
		}
	}()

	return stop, waitingChan
}

func (tracker *RenderTracker) wait(busy <-chan bool, waiting <-chan bool) {
	for {
		if !(<-busy || <-waiting) {
			return
		}
	}
}

func (tracker *RenderTracker) work(more []SharedRegionNumerics) {
	go func() {
		inst := RenderInput{
			Command: ThreadRender,
			Regions: more,
		}
		inputChan := <-tracker.schedule
		inputChan<- inst
	}()
}

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
	// Launch threads
	for i := uint32(0); i < tracker.config.Jobs; i++ {
		thread := tracker.factory.Build(tracker.ready[i], tracker.working[i], tracker.input[i], tracker.output[i])
		go thread.Run()
	}

	// Run the pool
	tracker.work([]SharedRegionNumerics{tracker.initialRegion})
	inputStop, busyChan := tracker.inputSchedule()
	outputStop, waitingChan := tracker.circulate()

	// Stop the pool
	tracker.wait(busyChan, waitingChan)
	inputStop<- true
	outputStop<- true

	// Shut down threads
	for _, input := range tracker.input {
		input <- RenderInput{Command: ThreadStop}
	}

	tracker.draw()
}