package main

import (
    "image"
    "image/png"
    "os"
    "io"
    "log"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

func main() {
    var input io.Reader = os.Stdin
    var output io.Writer = os.Stdout

    frch := lib.ReadInfoStream(input)
    imgch := make(chan image.Image)

    go func() {
        for frpkt := range frch {
            if frpkt.Err != nil {
                log.Fatal(frpkt.Err)
            }
            picture, renderErr := lib.Render(frpkt.Info)

            if renderErr != nil {
                log.Fatal("Render errror:", renderErr)
            }

            imgch<- picture
        }
        close(imgch)
    }()

    for picture := range imgch {
        encodeErr := png.Encode(output, picture)

        if encodeErr != nil {
            log.Fatal("Encoding error:", encodeErr)
        }
    }

}