package main

import (
    "image/png"
    "os"
    "io"
    "log"
    "runtime"
    "functorama.com/demo/libgodelbrot"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    var input io.Reader = os.Stdin
    var output io.Writer = os.Stdout

    desc, readErr := libgodelbrot.ReadInfo(input)

    if readErr != nil {
        log.Fatal("Error reading info: ", readErr)
    }

    picture, renderErr := libgodelbrot.Render(desc)

    if renderErr != nil {
        log.Fatal("Render errror:", renderErr)
    }

    encodeErr := png.Encode(output, picture)

    if encodeErr != nil {
        log.Fatal("Encoding error:", encodeErr)
    }
}