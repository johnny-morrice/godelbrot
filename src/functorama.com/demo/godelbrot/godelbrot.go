package main

import (
	"fmt"
	"flag"
	"image/png"
	"image"
	"log"
	"os"
	"runtime"
	"functorama.com/demo/libgodelbrot"
)

// Golang entry point
func main() {
		// Set number of cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	args := parseArguments()

	// Render the Mandelbrot set
	picture, renderError := renderFractal(args
	if renderError != nil {
		log.Fatal(renderError)
	}

	// Save the Mandelbrot set image to a file
	fileError := writePNGFile(args, picture)
	if fileError != nil {
		log.Fatal(renderError)
	}
}

// Structure representing our command line arguments
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
	storedPalette  string
	fixAspect 	   bool
	numericalSystem string
	glitchSamples uint
}

// Parse command line arguments into a `commandLine' structure
func parseArguments() commandLine {
	args := commandLine{}
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
	flag.UintVar(&args.glitchSamples, "regionGlitchSamples",
		libgodelbrot.DefaultRegionGlitchSampleSize, "Size of region render glitch-correncting sample set")
	flag.StringVar(&args.storedPalette, "storedPalette",
		"pretty", "Name of stored palette (pretty|redscale)")
	flag.StringVAr(&args.numericalSystem, "numerics",
		"auto", "Numerical system (auto|native|bigfloat)")
	flag.BoolVar(&args.fixAspect, "fixAspect",
		true, "Resize plane window to fit image aspect ratio")
	flag.Parse()

	return args
}

// Given the command line arguments, render the mandelbrot set
func renderFractal(args commandLine) (image.NRGBA, error) {
	renderDescription, validationError := extractRenderParameters(args)
	if validationError != nil {
		return nil, validationError
	}

	picture, renderError := libgodelbrot.Godelbrot(renderDescription)
	if renderError != nil {
		return nil, renderError
	}

	return picture, nil
}

// Validate and extract a render description from the command line arguments
func extractRenderParameters(args commandLine) (libgodelbrot.RenderDescription, error) {
	if args.iterateLimit > 255 {
		return nil, fmt.Errorf("iterateLimit out of bounds.  Valid values are: (0-255)")
	}

	if args.divergeLimit <= 0.0 {
		return nil, fmt.Errorf("divergeLimit out of bounds.  Valid values are: (> 0)")
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

	return description, nil
}

// Given the command line arguments and a picture, write the picture to a PNG
// file
func writePNGFile(args commandLine, picture image.NRGBA)
	file, fileError := os.Create(args.filename)

	if fileError != nil {
		log.Fatal(fileError)
	}
	defer file.Close()

	writeError := png.Encode(file, picture)

	if writeError != nil {
		log.Fatal(writeError)
	}
}
