package main

import (
	"errors"
	"flag"
	"functorama.com/demo/libgodelbrot"
	"image/png"
	"log"
	"os"
	"runtime"
)

type commandLine struct {
	iterateLimit   uint
	divergeLimit   float64
	width          uint
	height         uint
	filename       string
	realMin        float64
	realMax        float64
	imagMin        float64
	imagMax        float64
	zoom           float64
	mode           string
	frame          string
	regionCollapse uint
	renderThreads  uint
	threadBuffer   uint
	storedPalette  string
	fixAspect 	   bool
}

func parseArguments(args *commandLine) {
	realMin := real(libgodelbrot.MagicOffset)
	imagMax := imag(libgodelbrot.MagicOffset)
	realMax := realMin + real(libgodelbrot.MagicSetSize)
	imagMin := imagMax - imag(libgodelbrot.MagicSetSize)

	var renderThreads uint
	if cpus := runtime.NumCPU(); cpus > 1 {
		renderThreads = uint(cpus - 1)
	} else {
		renderThreads = 1
	}

	flag.UintVar(&args.iterateLimit, "iterateLimit", uint(libgodelbrot.DefaultIterations), "Maximum number of iterations")
	flag.Float64Var(&args.divergeLimit, "divergeLimit", libgodelbrot.DefaultDivergeLimit, "Limit where function is said to diverge to infinity")
	flag.UintVar(&args.width, "imageWidth", libgodelbrot.DefaultImageWidth, "Width of output PNG")
	flag.UintVar(&args.height, "imageHeight", libgodelbrot.DefaultImageHeight, "Height of output PNG")
	flag.StringVar(&args.filename, "output", "mandelbrot.png", "Name of output PNG")
	flag.Float64Var(&args.realMin, "realMin", realMin, "Leftmost position of complex plane projected onto PNG image")
	flag.Float64Var(&args.imagMax, "imagMax", imagMax, "Topmost position of complex plane projected onto PNG image")
	flag.Float64Var(&args.zoom, "zoom", libgodelbrot.DefaultZoom, "Zoom format")
	flag.Float64Var(&args.realMax, "realMax", realMax, "Rightmost position of complex plane projection")
	flag.Float64Var(&args.imagMin, "imagMin", imagMin, "Bottommost position of complex plane projection")
	flag.StringVar(&args.mode, "mode", "region", "Render mode.  Either 'sequence' or 'region'")
	flag.StringVar(&args.frame, "frame", "corner", "Coordinate frame.  Either 'corner' or 'zoom'")
	flag.UintVar(&args.regionCollapse, "collapse", libgodelbrot.DefaultCollapse, "Pixel width of region at which sequential render is forced")
	flag.UintVar(&args.renderThreads, "jobs", renderThreads, "Number of rendering threads in concurrent renderer")
	flag.UintVar(&args.threadBuffer, "buffer", libgodelbrot.DefaultBufferSize, "Size of per-thread buffer")
	flag.StringVar(&args.storedPalette, "storedPalette", "pretty", "Name of stored palette (pretty|redscale)")
	flag.BoolVar(&args.fixAspect, "fixAspect", true, "Fix aspect ratio: do not deform image")
	flag.Parse()
}

func extractRenderParameters(args commandLine) (*libgodelbrot.RenderConfig, error) {
	if args.iterateLimit > 255 {
		return nil, errors.New("iterateLimit out of bounds (uint8)")
	}

	if args.divergeLimit <= 0.0 {
		return nil, errors.New("divergeLimit out of bounds (positive float64)")
	}

	if args.zoom <= 0.0 {
		return nil, errors.New("zoom out of bounds (positive float64)")
	}

	parameters := libgodelbrot.RenderParameters{
		IterateLimit:   uint8(args.iterateLimit),
		DivergeLimit:   args.divergeLimit,
		Width:          args.width,
		Height:         args.height,
		TopLeft:        complex(args.realMin, args.imagMax),
		BottomRight:    complex(args.realMax, args.imagMin),
		Zoom:           args.zoom,
		RegionCollapse: args.regionCollapse,
		RenderThreads:  args.renderThreads,
		BufferSize:     args.threadBuffer,
		FixAspect:		args.fixAspect,
	}
	return parameters.Configure(), nil
}

func main() {
	args := commandLine{}
	parseArguments(&args)

	modes := map[string]libgodelbrot.Renderer{
		"sequence":   libgodelbrot.SequentialRender,
		"region":     libgodelbrot.RegionRender,
		"concurrent": libgodelbrot.ConcurrentRegionRender,
	}
	renderer := modes[args.mode]

	if renderer == nil {
		log.Fatal("Unknown renderer")
	}

	config, validationError := extractRenderParameters(args)
	if validationError != nil {
		log.Fatal(validationError)
	}

	palettes := map[string]libgodelbrot.PaletteFactory{
		"pretty": libgodelbrot.NewPrettyPalette,
		"redscale": libgodelbrot.NewRedscalePalette,
	}

	palette := palettes[args.storedPalette](config.IterateLimit)

	if palette == nil {
		log.Fatal("Unknown palette")
	}

	// Set number of cores
	runtime.GOMAXPROCS(int(config.RenderThreads) + 1)

	image, renderError := renderer(config, palette)
	if renderError != nil {
		log.Fatal(renderError)
	}

	file, fileError := os.Create(args.filename)

	if fileError != nil {
		log.Fatal(fileError)
	}
	defer file.Close()

	writeError := png.Encode(file, image)

	if writeError != nil {
		log.Fatal(writeError)
	}
}
