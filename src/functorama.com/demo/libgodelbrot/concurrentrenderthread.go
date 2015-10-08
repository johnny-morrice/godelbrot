type renderHeaps struct {
    escapePointHeap  *EscapePointHeap
    renderConfigHeap *RenderConfigHeap
}


// How much memory should we allocate in advance?
func newRenderOutput(config *RenderConfig) renderOutput {
    buffSize := config.BufferSize
    memberSize := config.RegionCollapse * config.RegionCollapse * config.BufferSize
    return renderOutput{
        uniformRegions: make([]RegionRenderContext, 0, buffSize),
        // What is the correct size for this?
        members: make([]pixelMember, 0, memberSize),
        // The magic number of 4 represents the splitting factor
        children: make([]RegionRenderContext, 0, buffSize),
    }
}

// A pass through the region rendering process, comprising many steps
func RegionRenderPass(config *RenderConfig, heaps renderHeaps, regions []RegionRenderContext) renderOutput {
    output := newRenderOutput(config)
    for _, region := range regions {
        RegionRenderStep(config, heaps, region, &output)
    }
    return output
}

func RegionRenderStep(config *RenderConfig, heaps renderHeaps, region RegionRenderContext, output *renderOutput) {
    if region.Collapse(config) {
        smallConfig := heaps.renderConfigHeap.Subconfig(region)
        MandelbrotSequence(smallConfig, func(i int, j int, member MandelbrotMember) {
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

// Implements a single render thread
func RegionRenderThead(threadNum uint, config *RenderConfig, inputChan <-chan renderInput, outputChan chan<- renderOutput) {
    heaps := renderHeaps{
        renderConfigHeap: NewRenderConfigHeap(tracker.context.Config, Kilo),
        escapePointHeap:  NewEscapePointHeap(K64),
    }
    for {
        input := <-inputChan
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