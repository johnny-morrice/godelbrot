package main

import (
    "flag"
    "io"
    "log"
    "os"
    "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

func main() {
    var input io.Reader = os.Stdin
    var output io.Writer = os.Stdout

    args := readArgs()

    z, zerr := libgodelbrot.ReadZoom(input)
    if zerr != nil {
        log.Fatal("Could not read zoom:", zerr)
    }

    frames, moverr := z.Movie(args.count)
    if moverr != nil {
        log.Fatal("Error zooming:", moverr)
    }

    for _, info := range frames {
        outerr := libgodelbrot.WriteInfo(output, info)
        if outerr != nil {
            log.Println("Error writing output:", outerr)
            break
        }
    }
}

func readArgs() params {
    args := params{}
    flag.UintVar(&args.count, "frames", 1, "Number of frames in zoom")
    flag.Parse()

    return args
}

type params struct {
    count uint
}