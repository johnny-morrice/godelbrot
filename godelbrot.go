package main

import (
    "fmt"
    "io"
    "os"
    "os/exec"
)

func main() {
    config := exec.Command("configbrot", os.Args[1:]...)
    render := exec.Command("renderbrot")

    confoutp, conferrp, conferr := pipes(config)
    if conferr != nil {
        fatal(conferr)
    }

    rendoutp, renderrp, renderr := pipes(render)
    if renderr != nil {
        fatal(renderr)
    }

    render.Stdin = confoutp

    tasks := []*exec.Cmd{config, render}
    for _, t := range tasks {
        err := t.Start()
        if err != nil {
            fatal(err)
        }
    }

    _, outerr := io.Copy(os.Stdout, rendoutp)
    if outerr != nil {
        fatal(outerr)
    }

    cerrcount, confcpyerr := io.Copy(os.Stderr, conferrp)
    if confcpyerr != nil {
        fatal(confcpyerr)
    }
    // If we read an error from configbrot, don't read an error from renderbrot
    if cerrcount == 0 {
        _, rndcpyerr := io.Copy(os.Stderr, renderrp)
        if rndcpyerr != nil {
            fatal(rndcpyerr)
        }
    }

    // Order of tasks is important!
    for _, t := range tasks {
        err := t.Wait()
        if err != nil {
            fmt.Fprintf(os.Stderr, "%v: %v\n", t.Path, err)
            os.Exit(2) // Different exit code for subprocess failure
        }
    }
}

func pipes(task *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
    outp, outerr := task.StdoutPipe()
    if outerr != nil {
        return nil, nil, outerr
    }

    errp, errerr := task.StderrPipe()
    if errerr != nil {
        return nil, nil, errerr
    }

    return outp, errp, nil
}

func fatal(err error) {
    fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
    os.Exit(1)
}