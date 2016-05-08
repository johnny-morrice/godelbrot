    package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"

    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
    "github.com/johnny-morrice/godelbrot/rest"
)

type params struct {
    jobs uint
    port uint
    bind string
    debug bool
    origins string
}

func main() {
    args := readArgs()

    info, readerr := lib.ReadInfo(os.Stdin)
    if readerr != nil {
        log.Printf("Info read error: %v", readerr)
    }

    const prefix = "/"
    h := rest.MakeWebservice(info, args.jobs, prefix)
    
    // Add CORS to the handler
    if args.origins != "" {
        h = newcorshandler(h, args.origins)
    }

    // Add logging to the handler
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

type corshandler struct {
    handler http.Handler
    origins map[string]bool
}

func newcorshandler(h http.Handler, origins string) http.Handler {
    ch := corshandler{}
    ch.handler = h

    ch.origins = map[string]bool{}
    for _, o := range strings.Split(origins, ",") {
        ch.origins[o] = true
    }

    return ch
}

func (ch corshandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    host := req.Host
    if ch.origins[host] {
        w.Header()["Access-Control-Allow-Origin"] = []string{host}
        log.Printf("Added CORS header for %v", host)
    }
    ch.handler.ServeHTTP(w, req)
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
    flag.StringVar(&args.origins, "origins", "", "Comma separated CORS client origins")
    flag.BoolVar(&args.debug, "debug", false, "Verbose logging")
    flag.Parse()
    return args
}