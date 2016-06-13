package main

import (
    "fmt"
    "flag"
    "log"
    "strconv"
    "os"
    "github.com/johnny-morrice/godelbrot/config"
    "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

func main() {
    output := os.Stdout

    args := parseArguments()

    req, inerr := newRequest(args)
    if inerr != nil {
        log.Fatal("Error forming request: ", inerr)
    }

    desc, gerr := libgodelbrot.Configure(req)

    if gerr != nil {
        log.Fatal("Error configuring Info: ", gerr)
    }

    outerr := libgodelbrot.WriteInfo(output, desc)
    if outerr != nil {
        log.Fatal("Error writing Info: ", outerr)
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
    fixAspect      string
    numerics string
    glitchSamples uint
    precision uint
    reconfigure bool
    palette string
}

// Parse command line arguments into a `commandLine' structure
func parseArguments() commandLine {
    args := commandLine{}

    bnds := []float64{
        real(libgodelbrot.MandelbrotMin),
        imag(libgodelbrot.MandelbrotMin),
        real(libgodelbrot.MandelbrotMax),
        imag(libgodelbrot.MandelbrotMax),
    }
    argbnds := make([]string, len(bnds))
    for i, c := range bnds {
        argbnds[i] = strconv.FormatFloat(c, 'e', -1, 64)
    }

    flag.UintVar(&args.iterateLimit, "iterlim",
        uint(libgodelbrot.DefaultIterations), "Maximum number of iterations")
    flag.Float64Var(&args.divergeLimit, "divlim",
        libgodelbrot.DefaultDivergeLimit, "Limit where function is said to diverge to infinity")
    flag.UintVar(&args.width, "width",
        libgodelbrot.DefaultImageWidth, "Width of output PNG")
    flag.UintVar(&args.height, "height",
        libgodelbrot.DefaultImageHeight, "Height of output PNG")
    flag.StringVar(&args.realMin, "rmin",
        argbnds[0], "Leftmost position on complex plane")
    flag.StringVar(&args.imagMin, "imin",
        argbnds[1], "Bottommost position on complex plane")
    flag.StringVar(&args.realMax, "rmax",
        argbnds[2], "Rightmost position on complex plane")
    flag.StringVar(&args.imagMax, "imax",
        argbnds[3], "Topmost position on complex plane")
    flag.StringVar(&args.mode, "render", "auto",
        "Render mode.  (auto|sequence|region)")
    flag.UintVar(&args.regionCollapse, "collapse",
        libgodelbrot.DefaultCollapse, "Pixel width of region at which sequential render is forced")
    flag.UintVar(&args.glitchSamples, "samples",
        libgodelbrot.DefaultRegionSamples, "Size of region sample set")
    flag.UintVar(&args.precision, "prec",
        libgodelbrot.DefaultPrecision, "Precision for big.Float render mode")
    flag.StringVar(&args.numerics, "numerics",
        "auto", "Numerical system (auto|native|bigfloat)")
    flag.StringVar(&args.palette, "palette", "grayscale", "(redscale|grayscale|pretty)")
    flag.StringVar(&args.fixAspect, "fix", "shrink", "Aspect ratio conservation (stretch|shrink|grow)")
    flag.BoolVar(&args.reconfigure, "reconf", false,
        "Reconfigure the render spec sent to stdin")
    flag.Parse()

    return args
}

// Validate and extract a render description from the command line arguments
func newRequest(args commandLine) (*config.Request, error) {
    user, uerr := userReq(args)

    if uerr != nil {
        return nil, uerr
    }

    var req *config.Request
    if args.reconfigure {
        desc, rerr := libgodelbrot.ReadInfo(os.Stdin)
        if rerr != nil {
            return nil, rerr
        }
        req = &desc.UserRequest
    } else {
        req = libgodelbrot.DefaultRequest()
    }

    argact := map[string]func(){
        "fix": func () {req.FixAspect = user.FixAspect},
        "palette": func () {req.PaletteCode = user.PaletteCode},
        "numerics": func () {req.Numerics = user.Numerics},
        "prec": func () {req.Precision = user.Precision},
        "jobs": func () {req.Jobs = user.Jobs},
        "collapse": func () {req.RegionCollapse = user.RegionCollapse},
        "render": func () {req.Renderer = user.Renderer},
        "iterlim": func () {req.IterateLimit = user.IterateLimit},
        "divlim": func () {req.DivergeLimit = user.DivergeLimit},
        "width": func () {req.ImageWidth = user.ImageWidth},
        "height": func () {req.ImageHeight = user.ImageHeight},
        "rmin": func () {req.RealMin = user.RealMin},
        "rmax": func () {req.RealMax = user.RealMax},
        "imin": func () {req.ImagMin = user.ImagMin},
        "imax": func () {req.ImagMax = user.ImagMax},
        "samples": func() {req.RegionSamples = user.RegionSamples},
        "reconf": func() {},
    }

    flag.Visit(func (fl *flag.Flag) {
        act := argact[fl.Name]
        if act == nil {
            log.Fatal("BUG: unknown action ", fl.Name)
        }
        act()
    })

    return req, nil
}

func userReq(args commandLine) (*config.Request, error) {
    const max8 = uint(^uint8(0))
    if args.iterateLimit > max8 {
        return nil, fmt.Errorf("iterateLimit out of bounds.  Valid values in range (0,%v)", max8)
    }

    if args.divergeLimit <= 0.0 {
        return nil, fmt.Errorf("divergeLimit out of bounds.  Valid values in range (0,)")
    }

    const max16 = uint(^uint16(0))
    if args.jobs > max16 {
        return nil, fmt.Errorf("jobs out of bounds.  Valid values in range (0,%v)", max16)
    }

    numerics := config.AutoDetectNumericsMode
    switch args.numerics {
    case "auto":
        // No change
    case "bigfloat":
        numerics = config.BigFloatNumericsMode
    case "native":
        numerics = config.NativeNumericsMode
    default:
        return nil, fmt.Errorf("Unknown numerics mode: %v", args.numerics)
    }

    renderer := config.AutoDetectRenderMode
    switch args.mode {
    case "auto":
        // No change
    case "sequence":
        renderer = config.SequenceRenderMode
    case "region":
        renderer = config.RegionRenderMode
    default:
        return nil, fmt.Errorf("Unknown render mode: %v", args.mode)
    }

    aspect := config.Shrink
    switch args.fixAspect {
    case "shrink":
        // No change
    case "stretch":
        aspect = config.Stretch
    case "grow":
        aspect = config.Grow
    default:
        return nil, fmt.Errorf("Unknown aspect fix strategy: %v", args.fixAspect)
    }

    req := &config.Request{}
    req.IterateLimit = uint8(args.iterateLimit)
    req.DivergeLimit = args.divergeLimit
    req.RealMin = args.realMin
    req.RealMax = args.realMax
    req.ImagMin = args.imagMin
    req.ImagMax = args.imagMax
    req.ImageWidth = args.width
    req.ImageHeight = args.height
    req.PaletteCode = args.palette
    req.FixAspect = aspect
    req.Renderer = renderer
    req.Numerics = numerics
    req.RegionCollapse = args.regionCollapse
    req.RegionSamples = args.glitchSamples
    req.Precision = args.precision

    return req, nil
}