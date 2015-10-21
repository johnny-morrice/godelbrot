package libgodelbrot

type RenderTracker struct {
	// nth element incremented when nth thread processing
	processing []uint32
	// input channels to threads
	input []chan renderInput
	// output channels to threads
	output []chan renderOutput
	// round robin input scheduler
	nextThread int
	// Local buffer of unprocessed regions
	buffer []SharedRegionNumerics
	// Concurrent render config
	config ConcurrentRegionParameters
	// drawing context for drawing onto image
	draw drawingContext
	// Processed uniform regions
	uniform []SharedRegionNumerics
	// Processed Mandelbrot points
	points []PixelMember
}

func NewRenderTracker(app GodelbrotApplication) *RenderTracker {
	config := app.ConcurrentConfig()
	tracker := RenderTracker{
		processing: make([]uint, config.Jobs),
		input:      make([]chan renderInput, config.Jobs),
		output:     make([]chan renderOutput, config.Jobs),
		buffer:     newBuffer(config.BufferSize),
		nextThread: 0,
		config:     config,
		draw:       app.DrawingContext(),
		uniform:    make([]SharedRegionNumerics, allocMedium),
		points:     make([]SharedRegionNumerics, allocMedium),
	}
	for i := 0; i < int(config.Jobs); i++ {
		tracker.processing[i] = 0
		inputChan := make(chan renderInput, allocSmall)
		outputChan := make(chan renderOutput, allocSmall)
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

	input := renderInput{command: render}

	// If we're busy, wait for buffer to fill before sending
	if tracker.busy() {
		if len(tracker.buffer) >= int(tracker.config.BufferSize) {
			input.regions = tracker.buffer
			tracker.cleanBuffer()
			tracker.sendInput(input)
		}
	} else {
		input.regions = tracker.buffer
		tracker.cleanBuffer()
		tracker.sendInput(input)
	}
}

// Send input and mark as busy
func (tracker *RenderTracker) sendInput(input renderInput) {
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
func (tracker *RenderTracker) step(output renderOutput) {

	// Give more work to the threads
	tracker.renderRegions(output.children)

	// Stash the completed areas
	tracker.uniform = append(tracker.uniform, output.uniformRegions)
	tracker.points = append(tracker.points, output.members)
}

// draw to the image
func (tracker *RenderTracker) draw() {
	// We must still ClaimExtrinsics to ensure we work with the local data in a cached object.

	for _, uniform := range tracker.uniform {
		uniform.ClaimExtrinsics()
		uniform.drawUniform()
	}

	// We do not need to claim any extrinsics here, because we are merely drawing a render result
	// that requires no extra context
	for _, member := range tracker.points {
		tracker.draw.drawPointAt(member.i, member.j, member.MandelbrotMember)
	}
}

func (tracker *RenderTracker) cleanBuffer() {
	tracker.buffer = newBuffer(tracker.context.Config.BufferSize)
}

func newBuffer(bufferSize uint) []SharedRegionNumerics {
	return make([]SharedRegionNumerics, 0, bufferSize)
}

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
	initialRegion := WholeRegion(tracker.context.Config)

	for i, inputChan := range tracker.input {
		outputChan := tracker.output[i]
		go RegionRenderProcess(uint(i), tracker.context.Config, inputChan, outputChan)
	}

	tracker.renderRegions([]SharedRegionNumerics{initialRegion})

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
		input <- renderInput{command: stop}
	}

	tracker.draw()
}
