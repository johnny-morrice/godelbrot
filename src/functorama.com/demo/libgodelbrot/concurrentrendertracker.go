package libgodelbrot

type RenderTracker struct {
    // nth element incremented when nth thread processing
    processing []uint
    // input channels to threads
    input []chan renderInput
    // output channels to threads
    output []chan renderOutput
    // round robin input scheduler
    nextThread int
    // Local buffer of unprocessed regions
    buffer []RegionRenderContext
    // Concurrent render config
    config ConcurrentRegionParameters
    // Drawing context for drawing onto image
    draw DrawingContext
}

func NewRenderTracker(app GodelbrotApplication) *RenderTracker {
    config := app.ConcurrentConfig()
    tracker := RenderTracker{
        processing: make([]uint, config.Jobs),
        input:      make([]chan renderInput, config.Jobs),
        output:     make([]chan renderOutput, config.Jobs),
        buffer:     newBuffer(config.BufferSize),
        nextThread: 0,
        config:    config,
        draw: app.DrawingContext(),
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
func (tracker *RenderTracker) Busy() bool {
    for _, wait := range tracker.processing {
        if wait > 0 {
            return true
        }
    }
    return false
}

// Render a set of regions using the next thread in the round robin scheme
func (tracker *RenderTracker) RenderRegions(regions []RegionRenderContext) {

    tracker.buffer = append(tracker.buffer, regions...)

    if len(tracker.buffer) == 0 {
        return
    }

    input := renderInput{command: render}

    // If we're busy, wait for buffer to fill before sending
    if tracker.Busy() {
        if len(tracker.buffer) >= int(tracker.config.BufferSize) {
            input.regions = tracker.buffer
            tracker.cleanBuffer()
            tracker.SendInput(input)
        }
    } else {
        input.regions = tracker.buffer
        tracker.cleanBuffer()
        tracker.SendInput(input)
    }
}

// Send input and mark as busy
func (tracker *RenderTracker) SendInput(input renderInput) {
    threadIndex := tracker.nextThread
    tracker.input[threadIndex] <- input
    tracker.processing[threadIndex]++

    // Rount robin
    tracker.nextThread++
    if tracker.nextThread >= len(tracker.processing) {
        tracker.nextThread = 0
    }
}

// Draw if the output is complete, otherwise hand it back in
func (tracker *RenderTracker) Draw(output renderOutput) {

    tracker.RenderRegions(output.children)

    for _, uniform := range output.uniformRegions {
        uniform.DrawUniform()
    }

    for _, member := range output.members {
        tracker.context.DrawPointAt(member.i, member.j, member.MandelbrotMember)
    }

}

func (tracker *RenderTracker) cleanBuffer() {
    tracker.buffer = newBuffer(tracker.context.Config.BufferSize)
}

func newBuffer(bufferSize uint) []RegionRenderContext {
    return make([]RegionRenderContext, 0, bufferSize)
}

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
    initialRegion := WholeRegion(tracker.context.Config)

    for i, inputChan := range tracker.input {
        outputChan := tracker.output[i]
        go RegionRenderProcess(uint(i), tracker.context.Config, inputChan, outputChan)
    }

    tracker.RenderRegions([]RegionRenderContext{initialRegion})

    for tracker.Busy() {
        for i, outputChan := range tracker.output {
            select {
            case output := <-outputChan:
                tracker.processing[i]--
                tracker.Draw(output)
            default:
                // Wait till next input
            }
        }
    }

    // Shut down threads
    for _, input := range tracker.input {
        input <- renderInput{command: stop}
    }
}