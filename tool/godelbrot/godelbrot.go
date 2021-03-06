package main

import (
	"fmt"
	"github.com/johnny-morrice/godelbrot/process"
	"os"
)

func main() {
	err := process.ConfigRender(os.Stdout, os.Stderr, os.Args[1:])
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
	os.Exit(1)
}
