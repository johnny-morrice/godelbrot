package main

import (
    "bytes"
    "flag"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "runtime/debug"
    "time"
    "github.com/johnny-morrice/godelbrot/rest"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

func main() {
    args := readArgs()

    web := newWebClient(args)

    if args.cycle || args.newrq {
        info, ierr := lib.ReadInfo(os.Stdin)
        fatalguard(ierr)
        web.req = info.GenRequest()
    }

    zoomparts := map[string]bool {
        "xmin": true,
        "xmax": true,
        "ymin": true,
        "ymax": true,
    }

    flag.Visit(func (fl *flag.Flag) {
        _, ok := zoomparts[fl.Name]
        if ok {
            web.zoom = true
        }
    })

    // Ugly
    var r io.Reader
    if args.newrq {
        rqi, err := web.newrq()
        fatalguard(err)
        reqr, jerr := jsonr(rqi)
        fatalguard(jerr)
        r = reqr
    } else if args.getrq != "" {
        rqi, err := web.getrq()
        fatalguard(err)
        reqr, jerr := jsonr(rqi)
        fatalguard(jerr)
        r = reqr
    } else if args.getimag != "" {
        png, err := web.getimag()
        fatalguard(err)
        r = png
    } else if args.cycle {
        png, err := web.cycle()
        fatalguard(err)
        r = png
    }

    _, cpyerr := io.Copy(os.Stdout, r)
    fatalguard(cpyerr)
}

type webclient struct {
    args params
    client http.Client
    req lib.Request
    zoom bool
    tick *time.Ticker

}

func newWebClient(args params) *webclient {
    web := &webclient{}
    web.args = args
    web.client.Timeout = time.Millisecond * time.Duration(web.args.timeout)
    return web
}

func (web *webclient) cycle() (io.Reader, error) {
    newresp, err := web.newrq()
    if err != nil {
        return nil, err
    }
    for {
        url := web.url(newresp.RQStatusURL)
        rqstat, err := web.getrqraw(url)
        if err != nil {
            return nil, err
        }
        switch rqstat.State {
        case "done":
            url := web.url(rqstat.ImageURL)
            return web.getimagraw(url)
        case "error":
            weberr := fmt.Errorf("RQGetResp error: %v", rqstat.Error)
            return nil, weberr
        case "wait":
            // NOP
        default:
            panic(fmt.Errorf("Unknown status: %v", rqstat.State))
        }
    }
}

func (web *webclient) newrq() (*rest.RQNewResp, error) {
    return web.newrqraw(web.url("renderqueue"))
}

func (web *webclient) newrqraw(url string) (*rest.RQNewResp, error) {
    renreq, rerr := web.renreq()
    if rerr != nil {
        return nil, addstack(rerr)
    }
    buff := &bytes.Buffer{}
    werr := rest.WriteReq(buff, renreq)
    if werr != nil {
        return nil, addstack(werr)
    }
    resp, err := web.post(url, "application/json", buff)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, httpError(resp)
    }
    defer resp.Body.Close()
    rqi := &rest.RQNewResp{}
    derr := web.decode(resp.Body, rqi)
    return rqi, addstack(derr)
}

func (web *webclient) getrq() (*rest.RQGetResp, error) {
    url := web.url(fmt.Sprintf("renderqueue/%v", web.args.getrq))
    return web.getrqraw(url)
}

func (web *webclient) getrqraw(url string) (*rest.RQGetResp, error) {
    resp, err := web.get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    rqi := &rest.RQGetResp{}
    derr := web.decode(resp.Body, rqi)
    return rqi, addstack(derr)
}

func (web *webclient) getimag() (io.Reader, error) {
    url := fmt.Sprintf("image/%v", web.args.getimag)
    return web.getimagraw(web.url(url))
}

func (web *webclient) getimagraw(url string) (io.Reader, error) {
    resp, err := web.get(url)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, httpError(resp)
    }
    defer resp.Body.Close()
    buff := &bytes.Buffer{}
    _, cpyerr := io.Copy(buff, resp.Body)
    return buff, addstack(cpyerr)
}

func (web *webclient) renreq() (*rest.RenderRequest, error) {
    renreq := &rest.RenderRequest{}
    renreq.Req = web.req
    if web.zoom {
        renreq.WantZoom = true
        renreq.Target.Xmin = web.args.xmin
        renreq.Target.Xmax = web.args.xmax
        renreq.Target.Ymin = web.args.ymin
        renreq.Target.Ymax = web.args.ymax
    }

    return renreq, nil
}

func (web *webclient) url(last string) string {
    args := web.args
    if web.args.prefix == "" {
        return fmt.Sprintf("http://%v:%v/%v/",
            args.addr, args.port, last)
    } else {
        return fmt.Sprintf("http://%v:%v/%v/%v",
                args.addr, args.port, args.prefix, last)
    }
}

func (web *webclient) get(url string) (r *http.Response, err error) {
    web.cautiously(func () {
        if web.args.debug {
            log.Printf("GET %v", url)
        }
        r, err = web.client.Get(url)
        if web.args.debug {
            web.reportResponse(r, err)
        }
    })
    return
}

func (web *webclient) post(url, ctype string, body io.Reader) (r *http.Response, err error) {
    web.cautiously(func () {
        if web.args.debug {
            log.Printf("POST %v", url)
        }
        r, err = web.client.Post(url, ctype, body)
        if web.args.debug {
            web.reportResponse(r, err)
        }
    })
    return
}

func (web *webclient) reportResponse(r *http.Response, err error) {
    if err != nil {
        log.Printf("Error: %v", err)
    }
    log.Printf("Status: %v", r.Status)
    ctypeHeads := r.Header["Content-Type"]
    if len(ctypeHeads) != 1 {
        log.Printf("Bad Content-Type header")
    } else {
        log.Printf("Content-Type: %v", ctypeHeads[0])
    }
}

func (web *webclient) cautiously(f func()) {
    if web.tick == nil {
        web.tick = time.NewTicker(time.Duration(web.args.ticktime) * time.Millisecond)
    } else {
        <-web.tick.C
    }
    f()
}

func (web *webclient) decode(r io.Reader, any interface{}) error {
    if web.args.debug {
        buff := &bytes.Buffer{}
        r = io.TeeReader(r, buff)
        derr := decode(r, any)
        log.Printf("Decoded: %v", buff.String())
        return derr
    }
    return decode(r, any)
}

func httpError(resp *http.Response) error {
    buff := &bytes.Buffer{}
    err := resp.Write(buff)
    if err != nil {
        panic(err)
    }
    return fmt.Errorf("Response:\n%v", buff)
}

func decode(r io.Reader, any interface{}) error {
    dec := json.NewDecoder(r)
    return dec.Decode(any)
}

func jsonr(any interface{}) (io.Reader, error) {
    buff := &bytes.Buffer{}
    enc := json.NewEncoder(buff)
    err := enc.Encode(any)
    return buff, err
}

func fatalguard(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v", err)
        os.Exit(1)
    }
}

func addstack(err error) error {
    if err == nil {
        return nil
    } else {
        return fmt.Errorf("%v\n%v", err, string(debug.Stack()))
    }

}


func readArgs() params {
    args := params{}
    flag.StringVar(&args.addr, "remote", "localhost", "Remote address of restfulbrot service")
    flag.UintVar(&args.port, "port", 9898, "Port of remote service")
    flag.StringVar(&args.prefix, "prefix", "", "Prefix of service URL")
    flag.UintVar(&args.timeout, "timeout", 1000, "Web request abort timeout (milliseconds)")
    flag.UintVar(&args.ticktime, "ticktime", 100, "Max one request per tick (milliseconds)")
    flag.BoolVar(&args.newrq, "newrq", false, "Add new item to render queue (info from stdin)")
    flag.StringVar(&args.getrq, "getrq", "", "Get status of render queue item")
    flag.StringVar(&args.getimag, "getimag", "", "Download fractal render (png to stdout)")
    flag.BoolVar(&args.cycle, "cycle", true,
        "Wait for fractal to render (info from stdin, png to stdout")
    flag.BoolVar(&args.debug, "debug", false, "Verbose debug mode")
    flag.Parse()
    return args
}

type params struct {
    addr string
    port uint
    prefix string
    timeout uint
    ticktime uint
    debug bool

    newrq bool
    getrq string
    getimag string
    cycle bool

    config string
    xmin uint
    xmax uint
    ymin uint
    ymax uint
}