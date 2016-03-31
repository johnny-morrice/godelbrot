package main

import (
    "image/png"
    "os"
    "io"
    "log"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

func main() {
    var input io.Reader = os.Stdin
    var output io.Writer = os.Stdout

    frames := []*lib.Info{}
    var frameErr error
    for {
        info, readerr := lib.ReadInfo(input)
        if readerr != nil {
            frameErr = readerr
            break
        }
        frames = append(frames, info)
    }

    framecnt := len(frames)
    if frameErr != nil {
        log.Printf("Error after %v frames: %v", framecnt, frameErr)
    }

    if framecnt == 0 {
        log.Fatal("No input frames found")
    }

    for _, info := range frames {
        picture, renderErr := lib.Render(info)

        if renderErr != nil {
            log.Fatal("Render errror:", renderErr)
        }

        encodeErr := png.Encode(output, picture)

        if encodeErr != nil {
            log.Fatal("Encoding error:", encodeErr)
        }
    }
}