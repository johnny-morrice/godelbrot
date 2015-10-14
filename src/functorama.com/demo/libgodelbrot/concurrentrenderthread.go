type RenderThread struct {
    app RenderApplication
    threadNum uint
    inputChan <-chan renderInput
    outputChan chan<- renderOutput
    config ConcurrentRenderParameters
}

// Easy method to create render threads
func RenderThreadFactory(app RenderApplication) {
    count := 0
    config := app.ConcurrentConfig()
    return func(inputChan <-chan renderInput, outputChan chan<- renderOutput) RenderThread {
        thread := RenderThread{
            threadNum: count,
            inputChan: inputChan,
            outputChan: outputChan,
            config: config,
        }
        count++
        return thread
    }
}

// Implements a single render thread
func (thread RenderThread) Run() {
    for {
        input := <-thread.inputChan
        switch input.command {
        case render:
            thread.outputChan <- thread.RegionRenderPass(input.regions)
        case stop:
            return
        default:
            panic(fmt.Sprintf("Unknown render command in thread %v: %v", threadNum, input.command))
        }
    }
}

// A pass through the region rendering process, comprising many steps
func (thread RenderThread) Pass(regions []RegionRenderContext) renderOutput {
    output := thread.createRenderOutput()
    for _, region := range regions {
        thread.Step(region, &output)
    }
    return output
}

// A single render step
func (thread RenderThread) Step(region RegionRenderContext, output *renderOutput) {
    if Collapse(region) {
        points := SequenceCollapse(region)
        output.members = append(output.members, points...)
        return
    }

    if Subdivide(region) {
        output.children = append(output.children, splitee.Children()...)
        return
    }

    output.uniformRegions = append(output.uniformRegions, region)
}

// Create a new output packet
func (thread RenderThread) createRenderOutput() renderOutput {
    buffSize := thread.config.BufferSize
    memberSize := thread.config.RegionCollapse * thread.config.RegionCollapse * thread.config.BufferSize
    return renderOutput{
        // Q: How much memory should we allocate in advance?
        uniformRegions: make([]RegionRenderContext, 0, buffSize),
        // What is the correct size for this?
        members: make([]pixelMember, 0, memberSize),
        children: make([]RegionRenderContext, 0, buffSize),
    }
}