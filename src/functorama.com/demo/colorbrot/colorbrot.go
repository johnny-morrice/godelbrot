package main

import (
    "image/png"
    "image"
    "os"
    "log"
    "flag"
    "errors"
    "fmt"
    "math"
    "functorama.com/demo/draw"
    "functorama.com/demo/base"
)

type commandLine struct {
    storedPalette  string
    iterateLimit uint
}

func parseCommand() *commandLine {
    args := &commandLine{}
    flag.StringVar(&args.storedPalette, "storedPalette",
        "pretty", "Name of stored palette (pretty|redscale)")
    flag.UintVar(&args.iterateLimit, "iterateLimit", 255, "Iterate limit [1, 255]")
    flag.Parse()
    return args
}

func createStoredPalette(code string, iterlim uint8) (draw.Palette, error) {
    palettes := map[string]draw.PaletteFactory{
        "redscale": draw.NewRedscalePalette,
        "pretty":   draw.NewPrettyPalette,
    }
    found := palettes[code]
    if found == nil {
        return nil, errors.New(fmt.Sprintf("Unknown palette: %v", code))
    }
    return found(iterlim), nil
}

func extractPalette(args *commandLine) (draw.Palette, error) {
    // Validate first
    if args.iterateLimit < 1 || args.iterateLimit > 255 {
        return nil, errors.New("iterateLimit must be between 1 and 255 inclusive")
    }

    return createStoredPalette(args.storedPalette, uint8(args.iterateLimit))
}

func main() {
    input := os.Stdin
    output := os.Stdout

    args:= parseCommand()
    palette, palErr := extractPalette(args)

    if palErr != nil {
        log.Fatal("Error extracting palette: ", palErr)
    }

    gray, decErr := png.Decode(input)

    if decErr != nil {
        log.Fatal("Error decoding PNG: ", decErr)
    }

    // CAUTION lossy conversion
    iterlim := uint8(args.iterateLimit)
    bnd := gray.Bounds()
    bright := image.NewNRGBA(bnd)
    scale := float64(0xff) / float64(0xffff)
    for x := bnd.Min.X; x < bnd.Max.X; x++ {
        for y := bnd.Min.Y; y < bnd.Max.Y; y++ {
            bigdiv, _, _, _ := gray.At(x, y).RGBA()
            invdiv := uint8(math.Floor(scale * float64(bigdiv)))
            member := base.BaseMandelbrot{
                InvDivergence: invdiv,
                InSet: invdiv == iterlim,
            }
            col := palette.Color(member)
            bright.Set(x, y, col)
        }
    }

    encErr := png.Encode(output, bright)

    if encErr != nil {
        log.Fatal("Error encoding PNG: ", encErr)
    }
}