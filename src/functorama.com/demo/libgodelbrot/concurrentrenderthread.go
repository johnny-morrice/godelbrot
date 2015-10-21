package libgodelbrot

type renderCommand uint

const (
	render = renderCommand(iota)
	stop
)

type renderInput struct {
	command renderCommand
	regions []SharedRegonNumerics
}

type renderOutput struct {
	uniformRegions []SharedRegonNumerics
	children       []SharedRegonNumerics
	members        []PixelMember
}

type renderThread struct {
	threadId   uint
	inputChan  <-chan renderInput
	outputChan chan<- renderOutput
	config     ConcurrentRenderParameters
}

type renderThreadFactory func(inputChan <-chan renderInput, outputChan chan<- renderOutput) RenderThread

// Easy method to create render threads
func newRenderThreadFactory(app RenderApplication) renderThreadFactory {
	count := 0
	config := app.ConcurrentConfig()
	return func(inputChan <-chan renderInput, outputChan chan<- renderOutput) RenderThread {
		thread := RenderThread{
			threadId:   count,
			inputChan:  inputChan,
			outputChan: outputChan,
			config:     config,
		}
		count++
		return thread
	}
}

// Implements a single render thread
func (thread RenderThread) run() {
	for {
		input := <-thread.inputChan
		switch input.command {
		case render:
			thread.outputChan <- thread.RegionRenderPass(input.regions)
		case stop:
			return
		default:
			panic(fmt.Sprintf("Unknown render command in thread %v: %v",
				thread.threadId, input.command))
		}
	}
}

// A pass through the region rendering process, comprising many steps
func (thread RenderThread) pass(regions []SharedRegionNumerics) renderOutput {
	output := thread.createRenderOutput()
	for _, region := range regions {
		thread.Step(region, &output)
	}
	return output
}

// A single render step
func (thread RenderThread) step(region SharedRegionNumerics, output *renderOutput) {
	// We are in a thread, so we must be sure that we have our own copy of the context
	region.GrabThreadPrototype(thread.threadId)
	// We use proxies to share objects, so we've got to ensure we're using the correct local data
	region.ClaimExtrinsics()

	collapse := thread.config.RegionCollapseSize

	if Collapse(region, collapse) {
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
		uniformRegions: make([]SharedRegionNumerics, 0, buffSize),
		// What is the correct size for this?
		members:  make([]pixelMember, 0, memberSize),
		children: make([]SharedRegionNumerics, 0, buffSize),
	}
}
