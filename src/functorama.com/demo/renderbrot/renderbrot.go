package main

import (
    "image/png"
    "os"
    "io"
    "bytes"
    "log"
    "functorama.com/demo/libgodelbrot"
)

func main() {
    var input io.Reader = os.Stdin
    var output io.Writer = os.Stdout

    buff := bytes.Buffer{}
    count, readErr := buff.ReadFrom(input)

    if readErr != nil {
        log.Fatal("Read error after ", count, "bytes:", readErr)
    }

    desc, jsonErr := libgodelbrot.FromJSON(buff.Bytes())

    if jsonErr != nil {
        log.Fatal("Error decoding JSON:", jsonErr)
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