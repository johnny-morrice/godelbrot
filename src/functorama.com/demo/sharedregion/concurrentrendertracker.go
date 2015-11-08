package sharedregion

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
	"functorama.com/demo/region"
)

type RenderTracker struct {
	// Number of jobs
	jobs uint16
	workers []*Worker
	workersDone chan bool
	stateChan chan workerState
	schedule chan chan<- RenderInput
	// Concurrent render config
	config SharedRegionConfig
	// drawing context for drawing onto image
	context draw.DrawingContext
	uniformChan chan SharedRegionNumerics
	memberChan chan base.PixelMember
	childChan chan SharedRegionNumerics
	// Thread factory
	factory *WorkerFactory
	// Initial region
	initialRegion SharedRegionNumerics
}

type workerState struct {
	workerId int
	waiting bool
}

func NewRenderTracker(app RenderApplication) *RenderTracker {
	config := app.SharedRegionConfig()
	iJobs := int(config.Jobs)
	factory := NewWorkerFactory(app)

	tracker := RenderTracker{
		jobs: config.Jobs,
		workers: make([]*Worker, config.Jobs),
		workersDone: make(chan bool),
		stateChan: make(chan workerState),
		config:     config,
		context:       app.DrawingContext(),
		initialRegion: app.SharedRegionFactory().Build(),
		uniformChan: make(chan SharedRegionNumerics),
		childChan: make(chan SharedRegionNumerics),
		memberChan: make(chan base.PixelMember),
	}

	outputChannels := RenderOutput{
		UniformRegions: tracker.uniformChan,
		Children: tracker.childChan,
		Members: tracker.memberChan,
	}
	for i := 0; i < iJobs; i++ {
		tracker.workers[i] = factory.Build(outputChannels)
	}

	return &tracker
}

// draw to the image
func (tracker *RenderTracker) drawUniforms() {
	for uniform := range tracker.uniformChan {
		uniform.GrabWorkerPrototype(tracker.jobs)
		uniform.ClaimExtrinsics()
		region.DrawUniform(tracker.context, uniform)
	}
}

func (tracker *RenderTracker) drawMembers() {
	for member := range tracker.memberChan {
		draw.DrawPoint(tracker.context, member)
	}
}

// We need to stop this
func (tracker *RenderTracker) circulate()  {
	for {
		select {
		case child := <-tracker.childChan:
			tracker.addWork(child)
			continue
		}

		select {
		case <-tracker.workersDone:
			return
		default:
			continue
		}
	}
}

func (tracker *RenderTracker) addWork(child SharedRegionNumerics) {
	input := RenderInput{
		Region: child,
	}
	// We need to feed back asynchronously
	// otherwise we will block the workers
	go func() {
		inputChan := <-tracker.schedule
		go func() {
			inputChan<- input
		}()
	}()
}

func (tracker *RenderTracker) scheduleWorkers() {

	for i, worker := range tracker.workers {
		go func() {
			for ready := range worker.WaitingChan {
				if ready {
					tracker.schedule<- worker.InputChan
				}
				tracker.stateChan<- workerState{i, ready}
			}
		}()
	}
}

func (tracker *RenderTracker) detectEnd() {
	workerWaiting := make([]bool, tracker.jobs)

	for state := range tracker.stateChan {
		workerWaiting[state.workerId] = state.waiting

		allWaiting := true
		for _, oneWait := range workerWaiting {
			if !oneWait {
				allWaiting = false
				break
			}
		}
		// Shut down the thread pool
		if allWaiting {
			tracker.workersDone<- true
		}
	}

}

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
	// Launch threads
	for _, worker := range tracker.workers {
		go worker.Run()
	}

	// Render fractal
	go func() { tracker.workers[0].InputChan<- RenderInput{tracker.initialRegion} }()
	go tracker.scheduleWorkers()
	go tracker.detectEnd()
	go tracker.drawUniforms()
	go tracker.drawMembers()

	// Circulate output to input until the fractal is drawn
	tracker.circulate()

	close(tracker.uniformChan)
	close(tracker.memberChan)

	// Shut down workers
	for _, worker := range tracker.workers {
		close(worker.InputChan)
	}
}