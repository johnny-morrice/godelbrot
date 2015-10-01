package libgodelbrot

import (
    "image"
    "fmt"
    "runtime"
)

type renderCommand uint
const (
    render = renderCommand(iota)
    stop = renderCommand(iota)
)

type renderInput struct {
    command renderCommand
    regions []Region
}

type renderResult uint

type pixelMember struct {
    i int
    j int
    MandelbrotMember
}

type renderOutput struct {
    uniformRegions []Region
    children []Region
    members []pixelMember
}

type RenderTracker struct {
    // nth element incremented when nth thread processing
    processing []uint
    // input channels to threads
    input [] chan renderInput
    // output channels to threads
    output [] chan renderOutput
    // round robin input scheduler
    nextThread int
    // Drawing context
    context DrawingContext
    // Local buffer of unprocessed regions
    buffer []Region
}

func NewRenderTracker(drawingContext DrawingContext) *RenderTracker {
    jobs := drawingContext.Config.RenderThreads
    tracker := RenderTracker{
        processing: make([]uint, jobs),
        input: make([] chan renderInput, jobs),
        output: make([] chan renderOutput, jobs),
        buffer: newBuffer(drawingContext.Config.BufferSize),
        nextThread: 0,
        context: drawingContext,
    }
    for i := 0; i < int(jobs); i++ {
        tracker.processing[i] = 0
        inputChan := make(chan renderInput, Meg)
        outputChan := make(chan renderOutput, Meg)
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
func (tracker *RenderTracker) RenderRegions(regions []Region) {

    tracker.buffer = append(tracker.buffer, regions...)

    if len(tracker.buffer) == 0 {
        return
    }

    input := renderInput{command: render}

    // If we're busy, wait for buffer to fill before sending
    if tracker.Busy() {
        if len(tracker.buffer) >= int(tracker.context.Config.BufferSize) {
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
        tracker.context.DrawUniform(uniform)
    }

    for _, member := range output.members {
        tracker.context.DrawPointAt(member.i, member.j, member.MandelbrotMember)
    }
    
}

func (tracker *RenderTracker) cleanBuffer() {
    tracker.buffer = newBuffer(tracker.context.Config.BufferSize)
}

func newBuffer(bufferSize uint) []Region {
    return make([]Region, 0, bufferSize)
}

type renderHeaps struct {
    escapePointHeap *EscapePointHeap
    renderConfigHeap *RenderConfigHeap
}

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
    initialRegion := WholeRegion(tracker.context.Config)

    for i, inputChan := range tracker.input {
        outputChan := tracker.output[i]
        heaps := renderHeaps{
            renderConfigHeap: NewRenderConfigHeap(tracker.context.Config, Meg), 
            escapePointHeap: NewEscapePointHeap(Meg),
        }
        go RegionRenderProcess(uint(i), heaps, tracker.context.Config, inputChan, outputChan)
    }

    tracker.RenderRegions([]Region{initialRegion})

    for tracker.Busy() {
        for i, outputChan := range tracker.output {
            select {
            case output := <- outputChan:
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

// How much memory should we allocate in advance?
func newRenderOutput(config *RenderConfig) renderOutput {
    buffSize := config.BufferSize
    memberSize := config.RegionCollapse * config.RegionCollapse * config.BufferSize
    return renderOutput{
        uniformRegions: make([]Region, 0, buffSize),
        // What is the correct size for this?
        members: make([]pixelMember, 0, memberSize),
        // The magic number of 4 represents the splitting factor
        children: make([]Region, 0, buffSize),
    }
}

// A pass through the region rendering process, comprising many steps
func RegionRenderPass(config *RenderConfig, heaps renderHeaps, regions []Region) renderOutput {
    output := newRenderOutput(config)
    for _, region := range regions {
        RegionRenderStep(config, heaps, region, &output)
    }
    return output
}

func RegionRenderStep(config *RenderConfig, heaps renderHeaps, region Region, output *renderOutput) {
    if region.Collapse(config) {
        smallConfig := heaps.renderConfigHeap.Subconfig(region)
        MandelbrotSequence(smallConfig, func (i int, j int, member MandelbrotMember) {
            output.members = append(output.members, pixelMember{i: i, j: j, MandelbrotMember: member})
        })
        return
    }

    subregion := region.Subdivide(config, heaps.escapePointHeap)
    if subregion.populated {
        output.children = append(output.children, subregion.children...)
        return
    }

    output.uniformRegions = append(output.uniformRegions, region)
}

// Implements a single render process
func RegionRenderProcess(threadNum uint, heaps renderHeaps, config *RenderConfig, inputChan <- chan renderInput, outputChan chan <- renderOutput) {
    for {
        input := <- inputChan
        switch input.command {
        case render:
            outputChan <- RegionRenderPass(config, heaps, input.regions)
        case stop:
            return
        default:
            panic(fmt.Sprintf("Unknown render command in thread %v: %v", threadNum, input.command))
        }
    }
}

func ConcurrentRegionRender(config *RenderConfig, palette Palette) (*image.NRGBA, error) {
    pic := config.BlankImage()
    ConcurrentRegionRenderImage(CreateContext(config, palette, pic))
    return pic, nil
}

func ConcurrentRegionRenderImage(drawingContext DrawingContext) {
    runtime.GOMAXPROCS(int(drawingContext.Config.RenderThreads) + 1)
    tracker := NewRenderTracker(drawingContext)
    tracker.Render()
}