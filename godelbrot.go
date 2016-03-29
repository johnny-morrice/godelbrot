package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "github.com/johnny-morrice/pipeline"
)

func main() {
    config := exec.Command("configbrot", os.Args[1:]...)
    render := exec.Command("renderbrot")

    pl := pipeline.New(&bytes.Buffer{}, os.Stdout, os.Stderr)
    pl.Chain(config, render)
    err := pl.Exec()
    if err != nil {
        fatal(err)
    }
}

func fatal(err error) {
    fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
    os.Exit(1)
}