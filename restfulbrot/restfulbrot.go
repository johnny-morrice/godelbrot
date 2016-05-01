package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
    "github.com/johnny-morrice/godelbrot/rest"
)

type params struct {
    jobs uint
    port uint
    bind string
    debug bool
}

func main() {
    args := readArgs()
    info, readerr := lib.ReadInfo(os.Stdin)
    if readerr != nil {
        log.Printf("Info read error: %v", readerr)
    }

    const prefix = "/"
    h := rest.MakeWebservice(info, args.jobs, prefix)
    if args.debug {
        log.Println("Running in debug mode")
        h = loghandler{handler: h}
        rest.Debug()
    }
    addr := fmt.Sprintf("%v:%v", args.bind, args.port)
    httperr := http.ListenAndServe(addr, h)
    if httperr != nil {
        log.Fatal(httperr)
    }
}

type loghandler struct {
    handler http.Handler
}

func (lh loghandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    log.Println(req.URL.RequestURI())
    lh.handler.ServeHTTP(w, req)
}

func readArgs() params {
    args := params{}
    flag.UintVar(&args.jobs, "jobs", 1, "Number of concurrent render jobs")
    flag.UintVar(&args.port, "port", 9898, "Port for webservice")
    flag.StringVar(&args.bind, "bind", "127.0.0.1", "Interface to bind against")
    flag.BoolVar(&args.debug, "debug", false, "Verbose logging")
    flag.Parse()
    return args
}