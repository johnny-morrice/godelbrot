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
	stateChan chan workerState
	processing []int
	workers []*Worker
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
	id uint16
	state int
}

type drawPacket struct {
	isRegion bool
	uniform SharedRegionNumerics
	points []base.PixelMember
}

func NewRenderTracker(app RenderApplication) *RenderTracker {
	output := RenderOutput{
		UniformRegions: make(chan WorkerRegionOut),
		Children: make(chan WorkerChildrenOut),
		Members: make(chan WorkerPixelOut),
	}

	config := app.SharedRegionConfig()
	workCount := config.Jobs - 1
	factory := NewWorkerFactory(app, output)

	tracker := RenderTracker{
		jobs: workCount,
		processing: make([]int, workCount),
		workers: make([]*Worker, workCount),
		stateChan: make(chan workerState),
		config:     config,
		context:       app.DrawingContext(),
		initialRegion: app.SharedRegionFactory().Build(),
		workerOutput: output,
	}

	for i := uint16(0); i < workCount; i++ {
		tracker.workers[i] = factory.Build()
	}

	return &tracker
}

func (tracker *RenderTracker) syncDrawing() <-chan drawPacket {
	// The number of goroutines we plan to spawn here
	const spawnCount = 2

	drawSync := make(chan drawPacket)
	wg := sync.WaitGroup{}
	wg.Add(spawnCount)

	go func() {
		for uni := range tracker.workerOutput.UniformRegions {
			tracker.stateChan<- workerState{id: uni.Id, state: -1}
			drawSync<- drawPacket{isRegion: true, uniform: uni.Region}
		}
		wg.Done()
	}()

	go func() {
		for detail := range tracker.workerOutput.Members {
			tracker.stateChan<- workerState{id: detail.Id, state: -1}
			drawSync<- drawPacket{isRegion: false, points: detail.Points}
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(drawSync)
	}()

	return drawSync
}

// draw to the image
func (tracker *RenderTracker) draw(pkt <-chan drawPacket) {
	for p := range pkt {
		if p.isRegion {
			p.uniform.GrabWorkerPrototype(tracker.jobs)
			p.uniform.ClaimExtrinsics()
			region.DrawUniform(tracker.context, p.uniform)
		} else {
			for _, px := range p.points {
				draw.DrawPoint(tracker.context, px)
			}
		}
	}
}

func (tracker *RenderTracker) circulate() {
	worker := uint16(0)
	for out := range tracker.workerOutput.Children {
		for i, child := range out.Children {
			// Order is important
			tracker.stateChan<- workerState{id: worker, state: 1}
			if i == 0 {
				tracker.stateChan<- workerState{id: out.Id, state: -1}
			}
			next := tracker.workers[worker]
			next.InputChan<- child

			// Round Robin
			worker++
			worker %= tracker.jobs
		}
	}
}

func (tracker *RenderTracker) wait() {
	for st := range tracker.stateChan {
		tracker.processing[st.id] += st.state

		over := true
		for _, p := range tracker.processing {
			if p > 0 {
				over = false
				break
			}
		}
		if over {
			return
		}
	}
}

func (tracker *RenderTracker) shutdown() {
	tracker.stopWorkers()
	tracker.workerOutput.Close()
	close(tracker.stateChan)
}

func (tracker *RenderTracker) stopWorkers() {
	for _, worker := range tracker.workers {
		close(worker.InputChan)
	}
	for _, worker := range tracker.workers {
		worker.Hold.Wait()
	}
}

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
	// Launch threads
	for _, worker := range tracker.workers {
		go worker.Run()
	}

	// Render fractal
	wid := 0
	tracker.processing[wid]++
	tracker.workers[wid].InputChan<- tracker.initialRegion
	go tracker.circulate()
	pkts := tracker.syncDrawing()
	go tracker.draw(pkts)

	tracker.wait()
	tracker.shutdown()
}