package process

import (
    "bytes"
    "io"
    "os/exec"
    "github.com/johnny-morrice/pipeline"
)

// Config creates a new Info, given the args, and sends it to stdout.
func Config(stdout, stderr io.Writer, args []string) error {
    config := configbrot(args)
    return runPipeCmd(config, &bytes.Buffer{}, stdout, stderr)
}

// Render sends a new fractal image to the passed stdout pipe, corresponding to the Info
// serialized in stdin.
func Render(stdin io.Reader, stdout, stderr io.Writer) error {
    render := renderbrot()
    return runPipeCmd(render, stdin, stdout, stderr)
}

// ConfigRender sends a new fractal image to the passed stdout pipe, corresponding to configbrot's
// processing of the args slice.
func ConfigRender(stdout, stderr io.Writer, args []string) error {
    config := configbrot(args)
    render := renderbrot()

    pl := pipeline.New(&bytes.Buffer{}, stdout, stderr)
    pl.Chain(config, render)
    return pl.Exec()
}

// Zoom magnifies a section of the Info read from stdin, sending it to stdout.
func Zoom(stdin io.Reader, stdout, stderr io.Writer, args[]string) error {
    zoom := zoombrot(args)
    return runPipeCmd(zoom, stdin, stdout, stderr)
}

// Zoom reads Info from stdin, and sends a fractal to stdout, returning the magnified Info,
// serialized as a io.Reader.
func ZoomRender(stdin io.Reader, stdout, stderr io.Writer, args []string) (io.Reader, error) {
    zoomBuff := &bytes.Buffer{}
    zoomerr := Zoom(stdin, zoomBuff, stderr, args)
    if zoomerr != nil {
        return nil, zoomerr
    }

    outbuff := &bytes.Buffer{}
    rendin := io.TeeReader(zoomBuff, outbuff)

    err := Render(rendin, stdout, stderr)

    return outbuff, err
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

func runPipeCmd(cmd *exec.Cmd, stdin io.Reader, stdout, stderr io.Writer) error {
    cmd.Stdin = stdin
    cmd.Stdout = stdout
    cmd.Stderr = stderr
    return cmd.Run()
}