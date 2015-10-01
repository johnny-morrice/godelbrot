package libgodelbrot

import (
    "image"
    "fmt"
)

type renderCommand uint
const (
    render = renderCommand(iota)
    stop = renderCommand(iota)
)

type renderInput struct {
    command renderCommand
    region *Region
}

type renderResult uint

const (
    uniform = renderResult(iota)
    small = renderResult(iota)
    divided = renderResult(iota)
)

type pixelMember struct {
    i int
    j int
    MandelbrotMember
}

type renderOutput struct {
    result renderResult
    uniformRegion *Region
    children []*Region
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
}

func NewRenderTracker(drawingContext DrawingContext) *RenderTracker {
    jobs := drawingContext.Config.RenderThreads
    tracker := RenderTracker{
        processing: make([]uint, jobs),
        input: make([] chan renderInput, jobs),
        output: make([] chan renderOutput, jobs),
        nextThread: 0,
        context: drawingContext,
    }
    for i := 0; i < int(jobs); i++ {
        tracker.processing[i] = 0
        inputChan := make(chan renderInput, 1024)
        outputChan := make(chan renderOutput, 1024)
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

// Render a region using the next thread in the round robin scheme
func (tracker *RenderTracker) RenderRegion(region *Region) {
    input := renderInput{
        command: render,
        region: region,
    }
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
    switch output.result {
    case divided:
        for _, child := range output.children {
            tracker.RenderRegion(child)
        }
    case uniform:
        tracker.context.DrawUniform(output.uniformRegion)
    case small:
        for _, member := range output.members {
            tracker.context.DrawPointAt(member.i, member.j, member.MandelbrotMember)
        }
    default:
        panic(fmt.Sprintf("Unknown output result: %v", output.result))
    }
}

// Render the Mandelbrot set concurrently
func (tracker *RenderTracker) Render() {
    initialRegion := WholeRegion(tracker.context.Config)

    for i, inputChan := range tracker.input {
        outputChan := tracker.output[i]
        go RegionRenderProcess(tracker.context.Config, inputChan, outputChan)
    }

    tracker.RenderRegion(initialRegion)

    for tracker.Busy() {
        for i, outputChan := range tracker.output {
            select {
            case firstOut := <- outputChan:
                // Wait on one less output from the thread
                tracker.processing[i]--
                tracker.Draw(firstOut)
            default:
                // Do nothing, wait till the next comes in
            }
        }
    }
}

// A single step through the region rendering process
func RegionRenderPass(config *RenderConfig, region *Region) renderOutput {
    if region.Collapse(config) {
        rect := region.Rect(config)
        area := rect.Dx() * rect.Dy()
        renderedPoints := renderOutput{
            result: small,
            members: make([]pixelMember, area, area),
        }
        index := 0
        MandelbrotSequence(config, func (i int, j int, member MandelbrotMember) {
            renderedPoints.members[index] = pixelMember{i: i, j: j, MandelbrotMember: member}
            index++
        })
        return renderedPoints
    }


    subregion := region.Subdivide(config)
    if subregion.populated {
        return renderOutput{
            result: divided,
            children: subregion.children,
        }
    }

    return renderOutput{
        result: uniform,
        uniformRegion: region,
    }
}

// Implements a single render process
func RegionRenderProcess(config *RenderConfig, inputChan <- chan renderInput, outputChan chan <- renderOutput) {
    input := <- inputChan
    switch input.command {
    case render:
        outputChan <- RegionRenderPass(config, input.region)
    case stop:
        return
    default:
        panic(fmt.Sprintf("Unknown render command: %v", input.command))
    }
}

func ConcurrentRegionRender(config *RenderConfig, palette Palette) (*image.NRGBA, error) {
    pic := config.BlankImage()
    ConcurrentRegionRenderImage(CreateContext(config, palette, pic))
    return pic, nil
}

func ConcurrentRegionRenderImage(drawingContext DrawingContext) {
    tracker := NewRenderTracker(drawingContext)
    tracker.Render()
}