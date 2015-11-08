package sharedregion

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
	"functorama.com/demo/region"
)

type RenderTracker struct {
	// Number of jobs
	jobs int
	workers []RenderThread
	schedule chan chan<- RenderInput
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
	factory := NewRenderThreadFactory(app)

	tracker := RenderTracker{
		jobs: int(config.Jobs),
		workers: make([]RenderThread, config.Jobs),
		config:     config,
		context:       app.DrawingContext(),
		memberChan: make(chan base.PixelMember),
		uniformChan: make(chan SharedRegionNumerics),
		initialRegion: app.SharedRegionFactory().Build(),
	}

	for i := 0; i < int(config.Jobs); i++ {
		tracker.workers[i] = factory.Build()
	}
	return &tracker
}

func (tracker *RenderTracker) stashUniforms(uniforms []SharedRegionNumerics) {
	for _, uniform := range uniforms {
		tracker.uniformChan<- uniform
	}
}

func (tracker *RenderTracker) stashMembers(points []base.PixelMember) {
	for _, member := range points {
		tracker.memberChan<- member
	}
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

func (tracker *RenderTracker) inputStep(busyChan chan<- bool) {
	busy := false
	for _, worker := range tracker.workers {
		select {
		case <-worker.ReadyChan:
			tracker.schedule<- worker.InputChan
		default:
			busy = true
			continue
		}
	}
	go func() { busyChan<- busy }()
}

func (tracker *RenderTracker) inputLoop() (chan<- bool, <-chan bool) {
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
			tracker.inputStep(busyChan)
		}
	}()

	return stop, busyChan
}

func (tracker *RenderTracker) circulateLoop() (chan<- bool, <-chan bool) {
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
			tracker.circulateStep(waitingChan)
		}
	}()

	return stop, waitingChan
}

// A single cycle of render output circulation
func (tracker *RenderTracker) circulateStep(waitingChan chan<- bool) {
	waiting := false
	for _, worker := range tracker.workers {
		select {
		case <-worker.WorkingChan:
			// Await result
			go func() {
				output := <-worker.OutputChan
				// Give more work to the threads
				go tracker.work(output.Children)
				// Stash stuff for the draw goroutine
				go tracker.stashUniforms(output.UniformRegions)
				go tracker.stashMembers(output.Members)
			}()
			waiting = true
		default:
			continue
		}
	}
	go func() {waitingChan<- waiting}()
}

// Wait until the workers are all waiting and not busy
func (tracker *RenderTracker) wait(busyChan <-chan bool, waitingChan <-chan bool) {
	for {
		waiting := <-waitingChan
		busy := <-busyChan
		if waiting && !busy {
			return
		}
	}
}

func (tracker *RenderTracker) work(more []SharedRegionNumerics) {
	inst := RenderInput{
		Command: ThreadRender,
		Regions: more,
	}
	inputChan := <-tracker.schedule
	inputChan<- inst
}

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
	// Launch threads
	for _, worker := range tracker.workers {
		go worker.Run()
	}

	// Render fractal
	go tracker.work([]SharedRegionNumerics{tracker.initialRegion})
	inputStop, busyChan := tracker.inputLoop()
	outputStop, waitingChan := tracker.circulateLoop()

	// Stop the pool
	tracker.wait(busyChan, waitingChan)
	inputStop<- true
	outputStop<- true

	// Shut down threads
	for _, worker := range tracker.workers {
		worker.InputChan<- RenderInput{Command: ThreadStop}
	}


}