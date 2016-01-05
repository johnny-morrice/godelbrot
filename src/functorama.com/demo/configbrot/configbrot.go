package main

import (
    "fmt"
    "flag"
    "log"
    "runtime"
    "strconv"
    "functorama.com/demo/libgodelbrot"
)

// Golang entry point
func main() {
    // Set number of cores
    runtime.GOMAXPROCS(runtime.NumCPU())

    args := parseArguments()
    request, argErr := extractRenderParameters(args)
    if argErr != nil {
        log.Fatal("Error:", argErr)
    }

    info, godelErr := libgodelbrot.AutoConf(request)

    if godelErr != nil {
        log.Fatal(godelErr)
    }

    text, jsonErr := libgodelbrot.ToJSON(info)
    if jsonErr == nil {
        fmt.Printf("%s", text)
    } else {
        log.Fatal("Error creating JSON:", jsonErr)
    }
}

// Structure representing our command line arguments
type commandLine struct {
    iterateLimit   uint
    divergeLimit   float64
    width          uint
    height         uint
    realMin        string
    realMax        string
    imagMin        string
    imagMax        string
    mode           string
    regionCollapse uint
    jobs  uint
    storedPalette  string
    fixAspect      bool
    numerics string
    glitchSamples uint
    precision uint
}

// Parse command line arguments into a `commandLine' structure
func parseArguments() commandLine {
    args := commandLine{}

    components := []float64{
        real(libgodelbrot.MandelbrotMin),
        imag(libgodelbrot.MandelbrotMin),
        real(libgodelbrot.MandelbrotMax),
        imag(libgodelbrot.MandelbrotMax),
    }
    bounds := make([]string, len(components))
    for i, num := range components {
        bounds[i] = strconv.FormatFloat(num, 'e', -1, 64)
    }


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
    flag.StringVar(&args.realMin, "realMin",
        bounds[0], "Leftmost position on complex plane")
    flag.StringVar(&args.imagMin, "imagMin",
        bounds[1], "Bottommost position on complex plane")
    flag.StringVar(&args.realMax, "realMax",
        bounds[2], "Rightmost position on complex plane")
    flag.StringVar(&args.imagMax, "imagMax",
        bounds[3], "Topmost position on complex plane")
    flag.StringVar(&args.mode, "mode", "auto",
        "Render mode.  (auto|sequence|region|concurrent)")
    flag.UintVar(&args.regionCollapse, "collapse",
        libgodelbrot.DefaultCollapse, "Pixel width of region at which sequential render is forced")
    flag.UintVar(&args.jobs, "jobs",
        renderThreads, "Number of rendering threads in concurrent renderer")
    flag.UintVar(&args.glitchSamples, "regionGlitchSamples",
        libgodelbrot.DefaultGlitchSamples, "Size of region render glitch-correncting sample set")
    flag.UintVar(&args.precision, "prec",
        libgodelbrot.DefaultPrecision, "Precision for big.Float render mode")
    flag.StringVar(&args.storedPalette, "storedPalette",
        "pretty", "Name of stored palette (pretty|redscale)")
    flag.StringVar(&args.numerics, "numerics",
        "auto", "Numerical system (auto|native|bigfloat)")
    flag.BoolVar(&args.fixAspect, "fixAspect",
        true, "Resize plane window to fit image aspect ratio")
    flag.Parse()

    return args
}

// Validate and extract a render description from the command line arguments
func extractRenderParameters(args commandLine) (*libgodelbrot.Request, error) {
    if args.iterateLimit > 255 {
        return nil, fmt.Errorf("iterateLimit out of bounds.  Valid values in range (0,255)")
    }

    if args.divergeLimit <= 0.0 {
        return nil, fmt.Errorf("divergeLimit out of bounds.  Valid values in range (0,)")
    }

    const max16 = uint(^uint16(0))
    if args.jobs > max16 {
        return nil, fmt.Errorf("jobs out of bounds.  Valid values in range (0, 65535)")
    }

    numerics := libgodelbrot.AutoDetectNumericsMode
    switch args.numerics {
    case "auto":
        // No change
    case "bigfloat":
        numerics = libgodelbrot.BigFloatNumericsMode
    case "native":
        numerics = libgodelbrot.NativeNumericsMode
    default:
        log.Fatal("Unknown numerics mode:", args.numerics)
    }

    renderer := libgodelbrot.AutoDetectRenderMode
    switch args.mode {
    case "auto":
        // No change
    case "sequence":
        renderer = libgodelbrot.SequenceRenderMode
    case "region":
        renderer = libgodelbrot.RegionRenderMode
    case "concurrent":
        renderer = libgodelbrot.SharedRegionRenderMode
    default:
        log.Fatal("Unknown render mode:", args.mode)
    }

    description := &libgodelbrot.Request {
        RealMin: args.realMin,
        RealMax: args.realMax,
        ImagMin: args.imagMin,
        ImagMax: args.imagMax,
        ImageWidth: args.width,
        ImageHeight: args.height,
        PaletteType: libgodelbrot.StoredPalette,
        PaletteCode: args.storedPalette,
        FixAspect: args.fixAspect,
        Numerics: numerics,
        Renderer: renderer,
        Jobs: uint16(args.jobs),
    }

    return description, nil
}