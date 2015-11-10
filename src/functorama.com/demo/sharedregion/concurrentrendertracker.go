package sharedregion

import (
	"sync"
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
	workerOutput RenderOutput
	// Thread factory
	factory *WorkerFactory
	// Initial region
	initialRegion SharedRegionNumerics
}

type workerState struct {
	workerId int
	waiting bool
}

type drawPacket struct {
	isRegion bool
	uniform SharedRegionNumerics
	point base.PixelMember
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
		workerOutput: RenderOutput{
			UniformRegions: make(chan SharedRegionNumerics),
			Children: make(chan SharedRegionNumerics),
			Members: make(chan base.PixelMember),
		}
	}

	for i := 0; i < iJobs; i++ {
		tracker.workers[i] = factory.Build(tracker.workerOutput)
	}

	return &tracker
}

func (tracker *RenderTracker) syncDrawing() chan<- drawPacket {
	// The number of goroutines we plan to spawn here
	const spawnCount = 2

	drawSync := make(chan drawPacket)
	wg := sync.WaitGroup{}
	wg.Add(spawnCount)

	go func() {
		for uni := range tracker.uniformChan {
			drawSync<- drawPacket{isRegion: true, uniform: uni}
		}
		wg.Done()
	}

	go func() {
		for member := range tracker.memberChan {
			drawSync<- drawPacket{isRegion: false, point: member}
		}
		wg.Done()
	}

	go func() {
		wg.Wait()
		close(drawSync)
	}

	return drawSync
}

// draw to the image
func (tracker *RenderTracker) draw(packets chan<- drawPacket) {
	for packet := range packets {
		if packet.isRegion {
			uniform.GrabWorkerPrototype(tracker.jobs)
			uniform.ClaimExtrinsics()
			region.DrawUniform(tracker.context, uniform)
		} else {
			draw.DrawPoint(tracker.context, member)
		}
	}
}

// We need to stop this
func (tracker *RenderTracker) circulate()  {
	shutdown := false
	for {
		select {
		case child := <-tracker.childChan:
			shutdown = false
			tracker.addWork(child)
			continue
		case <-tracker.workersDone:
			shutdown = true
		default:
			if shutdown {
				return
			}
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

func (tracker *RenderTracker) scheduleStep() chan<- bool {
	for running {
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

func (tracker *RenderTracker) scheduleWorkers() {
	return loop(tracker.scheduleStep)
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

		// Indicate that the workers have finished
		if allWaiting {
			// Wait on each worker
			// They may have data to send
			for _, worker := range tracker.workers {
				worker.hold.Wait()
			}
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
	go tracker.detectEnd()
	stop := tracker.scheduleWorkers()
	packets := tracker.syncDrawing(packets)
	go tracker.draw(packets)

	// Circulate output to input until the fractal is drawn
	tracker.circulate()

	close(tracker.uniformChan)
	close(tracker.memberChan)

	// Shut down workers
	for _, worker := range tracker.workers {
		close(worker.InputChan)
	}
}

func loop(f func()) chan<- bool {
	stop := make(chan bool)
	go func() { 
		for {
			// Check for shutdown
			select {
			case <-stop:
				return
			default:
			}
			f()
		}
	}
}