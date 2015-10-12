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
	realMin        string
	realMax        string
	imagMin        string
	imagMax        string
	mode           string
	regionCollapse uint
	renderThreads  uint
	threadBuffer   uint
	storedPalette  string
	fixAspect 	   bool
	numericalSystem string
	glitchSamples uint
}

func parseArguments(args *commandLine) {
	realMin := string(real(libgodelbrot.MagicOffset))
	imagMax := string(imag(libgodelbrot.MagicOffset))
	realMax := string(realMin + real(libgodelbrot.MagicSetSize))
	imagMin := string(imagMax - imag(libgodelbrot.MagicSetSize))

	var renderThreads uint
	if cpus := runtime.NumCPU(); cpus > 1 {
		renderThreads = uint(cpus - 1)
	} else {
		renderThreads = 1
	}

	flag.UintVar(&args.iterateLimit, "iterateLimit", 
		uint(libgodelbrot.DefaultIterations), "Maximum number of iterations")
	flag.Float64Var(&args.divergeLimit, "divergeLimit", 
		libgodelbrot.DefaultDivergeLimit, "Limit where function is said to diverge to infinity")
	flag.UintVar(&args.width, "imageWidth", 
		libgodelbrot.DefaultImageWidth, "Width of output PNG")
	flag.UintVar(&args.height, "imageHeight", 
		libgodelbrot.DefaultImageHeight, "Height of output PNG")
	flag.StringVar(&args.filename, "output", 
		"mandelbrot.png", "Name of output PNG")
	flag.StringVar(&args.realMin, "realMin", 
		realMin, "Leftmost position of complex plane projected onto PNG image")
	flag.StringVar(&args.imagMax, "imagMax", 
		imagMax, "Topmost position of complex plane projected onto PNG image")
	flag.StringVar(&args.realMax, "realMax", 
		realMax, "Rightmost position of complex plane projection")
	flag.StringVar(&args.imagMin, "imagMin", 
		imagMin, "Bottommost position of complex plane projection")
	flag.StringVar(&args.mode, "mode", "auto", 
		"Render mode.  (auto|sequence|region|concurrent)")
	flag.UintVar(&args.regionCollapse, "collapse", 
		libgodelbrot.DefaultCollapse, "Pixel width of region at which sequential render is forced")
	flag.UintVar(&args.renderThreads, "jobs", 
		renderThreads, "Number of rendering threads in concurrent renderer")
	flag.UintVar(&args.threadBuffer, "buffer", 
		libgodelbrot.DefaultBufferSize, "Size of per-thread buffer")
	flag.UintVar(&args.glitchSamples, "regionGlitchSamples",
		libgodelbrot.DefaultRegionGlitchSampleSize, "Size of region render glitch-correncting sample set")
	flag.StringVar(&args.storedPalette, "storedPalette", 
		"pretty", "Name of stored palette (pretty|redscale)")
	flag.StringVAr(&args.numericalSystem, "numerics",
		"auto", "Numerical system (auto|native|bigfloat)")
	flag.BoolVar(&args.fixAspect, "fixAspect", 
		true, "Resize plane window to fit image aspect ratio")
	flag.Parse()
}

func extractRenderParameters(args commandLine) (libgodelbrot.RenderContext, error) {
	if args.iterateLimit > 255 {
		return nil, errors.New("iterateLimit out of bounds (uint8)")
	}

	if args.divergeLimit <= 0.0 {
		return nil, errors.New("divergeLimit out of bounds (positive float64)")
	}

	description := libgodelbrot.RenderDescription {
		RealMin: args.realMin,
		RealMax: args.realMax,
		ImagMin: args.imagMin,
		ImagMax: args.imagMax,
		ImageWidth: args.imageWidth,
		ImageHeight: args.imageHeight,
		ThreadBufferSize: args.threadBuffer,
		PaletteType: libgodelbrot.StoredPalette,
		PaletteCode: args.storedPalette,
		FixAspect: args.fixAspect,
		Numerics: args.numerics,
		Renderer: mode,
		Jobs: args.jobs,
	}
	context := description.CreateInitialRenderContext()
	
	return context, nil
}

func main() {
		// Set number of cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	args := commandLine{}
	parseArguments(&args)

	context, validationError := extractRenderParameters(args)
	if validationError != nil {
		log.Fatal(validationError)
	}

	image, renderError := context.Render()
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
