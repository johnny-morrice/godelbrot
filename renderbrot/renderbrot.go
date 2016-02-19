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

    desc, readErr := lib.ReadInfo(input)

    if readErr != nil {
        log.Fatal("Error reading info: ", readErr)
    }

    picture, renderErr := lib.Render(desc)

    if renderErr != nil {
        log.Fatal("Render errror:", renderErr)
    }

    encodeErr := png.Encode(output, picture)

    if encodeErr != nil {
        log.Fatal("Encoding error:", encodeErr)
    }
}