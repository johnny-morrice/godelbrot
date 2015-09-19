package main

import (
    "flag"
    "log"
    "os"
    "image/png"
    "errors"
)

type commandLine struct {
    iterateLimit uint
    divergeLimit float64
    width uint
    height uint
    filename string
    xOffset float64
    yOffset float64
    zoom float64
    mode string   
}

func parseArguments(args *commandLine) {
    flag.UintVar(&args.iterateLimit, "iterateLimit", 255, "Maximum number of iterations")
    flag.Float64Var(&args.divergeLimit, "divergeLimit", 4.0 "Limit where function is said to diverge to infinity")
    flag.UintVar(&args.width, "imageWidth", 800, "Width of output PNG")
    flag.UintVar(&args.height, "imageHeight", 600, "Height of output PNG")
    flag.StringVar(&args.filename, "filename", "mandelbrot.png", "Name of output PNG")
    flag.Float64Var(&args.xOffset, "xOffset", -1.5, "Leftmost position of complex plane projected onto PNG image")
    flag.Float64Var(&args.yOffset, "yOffset", 1.0, "Topmost position of complex plane projected onto PNG image")
    flag.Float64Var(&args.zoom, "zoom", 1.0, "Look into the eyeball")
    flag.StringVar(&args.mode, "mode", "sequential", "Render mode")
}

func extractRenderParameters(args commandLine) (*RenderParameters, error) {
    if args.iterateLimit > 255 {
        return nil, errors.New("iterateLimit out of bounds (uint8)")
    }

    if args.divergeLimit <= 0.0 {
        return nil, errors.New("divergeLimit out of bounds (positive float64)")
    }

    if args.zoom <= 0.0 {
        return nil, errors.New("zoom out of bounds (positive float64)")
    }

    return &RenderParameters{
        IterateLimit: args.iterateLimit,
        DivergeLimit: args.divergeLimit,
        Width: args.width,
        Height: args.height,
        XOffset: args.xOffset,
        YOffset: args.yOffset,
        Zoom: args.zoom,
    }
}

func main() {
    commandLineArguments := commandLine{}
    parseArguments(&commandLine)

    renderer := nil
    switch mode {
    case "sequential":
        renderer := NewSequentialRenderer()
    default:
        log.Fatal("No renderer specified")
    }

    renderParameters, validationError := extractRenderParameters(commandLineArguments)
    if validationError != nil {
        log.Fatal(validationError)
    }

    image, renderError := renderer.Render(renderParameters)
    if renderError != nil {
        log.Fatal(renderError)
    }

    file, fileError := os.Create(commandLineArguments.filename)

    if fileError != nil {
        log.Fatal(fileError)
    }
    defer os.Close(file)

    writeError := png.Encode(file, image)

    if writeError != nil {
        log.Fatal(writeError)
    }
}