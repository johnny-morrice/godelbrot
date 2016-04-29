package rest

import (
    "bufio"
    "bytes"
    "io"
    "log"
    "github.com/johnny-morrice/godelbrot/process"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

type renderbuffers struct {
    png bytes.Buffer
    info bytes.Buffer
    nextinfo bytes.Buffer
    report bytes.Buffer
}

func (rb renderbuffers) logReport() {
    sc := bufio.NewScanner(&rb.report)
    for sc.Scan() {
        err := sc.Err()
        if err != nil {
            log.Printf("Error while printing error (omg!): %v", err)
        }
        log.Println(sc.Text())
    }
}

func (rb renderbuffers) input(info *lib.Info) error {
    return lib.WriteInfo(&rb.info, info)
}

// renderservice renders fractals
type renderservice struct {
    s sem
}

// makeRenderService creates a render service that allows at most `concurrent` concurrent tasks.
func makeRenderservice(concurrent uint) renderservice {
    rs := renderservice{}
    rs.s = semaphor(concurrent)
    return rs
}

// render a fractal into the renderbuffers
func (rs renderservice) render(rbuf renderbuffers, zoomArgs []string) error {
    rs.s.acquire(1)
    var err error
    if zoomArgs == nil || len(zoomArgs) == 0 {
        next, zerr := process.ZoomRender(&rbuf.info, &rbuf.png, &rbuf.report, zoomArgs)
        err = zerr
        if err == nil {
            _, err = io.Copy(&rbuf.nextinfo, next)
        }
    } else {
        tee := io.TeeReader(&rbuf.info, &rbuf.nextinfo)
        err = process.Render(tee, &rbuf.png, &rbuf.report)
    }
    rs.s.release(1)
    return err
}