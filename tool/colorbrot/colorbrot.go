package main

import (
    "log"
    "flag"
    "os"
    "image/png"
    "github.com/johnny-morrice/godelbrot"
)

type commandLine struct {
    config string
}

func parseCommand() *commandLine {
    args := &commandLine{}
    flag.StringVar(&args.config, "config", "(none)", "Path to config file")
    flag.Parse()
    return args
}

func readInfo(args *commandLine) (*godelbrot.Info, error) {
    if args.config == "(none)" {
        desc := &godelbrot.Info{}
        desc.PaletteType = godelbrot.Pretty
        return desc, nil
    } else {
        f, err := os.Open(args.config)
        if err != nil {
            return nil, err
        }
        defer f.Close()
        return godelbrot.ReadInfo(f)
    }
}

func main() {
    input := os.Stdin
    output := os.Stdout

    args:= parseCommand()
    desc, argErr := readInfo(args)

    if argErr != nil {
        log.Fatal("Error extracting palette: ", argErr)
    }

    gray, decErr := png.Decode(input)

    if decErr != nil {
        log.Fatal("Error decoding PNG: ", decErr)
    }

    bright := godelbrot.Recolor(desc, gray)

    encErr := png.Encode(output, bright)

    if encErr != nil {
        log.Fatal("Error encoding PNG: ", encErr)
    }
}