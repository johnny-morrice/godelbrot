package process

import (
    "bytes"
    "io"
    "os/exec"
    "github.com/johnny-morrice/pipeline"
)

// Render sends a new fractal image to the passed stdout pipe, corresponding to configbrot's
// processing of the args slice.
func Render(stdout, stderr io.Writer, args []string) error {
    config := configbrot(args)
    render := renderbrot()

    pl := pipeline.New(&bytes.Buffer{}, stdout, stderr)
    pl.Chain(config, render)
    return pl.Exec()
}

// Zoom reads *Info from stdin, and sends a fractal to stdout, returning the next *Info,
// serialized as a io.Reader.
func Zoom(stdin io.Reader, stdout, stderr io.Writer, args []string) (io.Reader, error) {
    zoom := zoombrot(args)
    zoom.Stdin = stdin
    zoom.Stderr = stderr

    zoombytes, zoomerr := zoom.Output()
    if zoomerr != nil {
        return nil, zoomerr
    }

    rendbuff, outbuff := bytes.NewBuffer(zoombytes), bytes.NewBuffer(zoombytes)
    render := renderbrot()
    render.Stdin = rendbuff
    render.Stdout = stdout
    render.Stderr = stderr

    err := render.Run()

    if err != nil {
        return nil, err
    }

    return outbuff, nil
}

func zoombrot(args []string) *exec.Cmd {
    return exec.Command("zoombrot", args...)
}

func configbrot(args []string) *exec.Cmd {
    return exec.Command("configbrot", args...)
}

func renderbrot() *exec.Cmd {
    return exec.Command("renderbrot")
}