package main

import (
    "flag"
    "io"
    "log"
    "os"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

func main() {
    var input io.Reader = os.Stdin
    var output io.Writer = os.Stdout

    args := readArgs()
    validerr := args.zt.Validate()
    if validerr != nil {
        log.Fatal(validerr)
    }

    info, readerr := lib.ReadInfo(input)
    if readerr != nil {
        log.Fatal("Could not read info:", readerr)
    }

    z := lib.Zoom{}
    z.ZoomTarget = args.zt
    z.Prev = *info

    frames, moverr := z.Movie()
    if moverr != nil {
        log.Fatal("Error zooming:", moverr)
    }

    for _, info := range frames {
        outerr := lib.WriteInfo(output, info)
        if outerr != nil {
            log.Println("Error writing output:", outerr)
            break
        }
    }
}

func readArgs() params {
    args := params{}
    flag.UintVar(&args.zt.Frames, "frames", 1, "Number of frames in zoom")
    flag.UintVar(&args.zt.Xmin, "xmin", 0, "X-Min")
    flag.UintVar(&args.zt.Xmax, "xmax", 0, "X-Max")
    flag.UintVar(&args.zt.Ymin, "ymin", 0, "Y-Min")
    flag.UintVar(&args.zt.Ymax, "ymax", 0, "Y-Max")
    flag.BoolVar(&args.zt.Reconfigure, "reconf", true, "Reconfigure magnified request")
    flag.BoolVar(&args.zt.UpPrec, "incprec", true, "Increase precision for zoom")
    flag.Parse()

    return args
}

type params struct {
    zt lib.ZoomTarget
}