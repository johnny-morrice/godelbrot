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
}

func main() {
    args := readArgs()
    info, readerr := lib.ReadInfo(os.Stdin)
    if readerr != nil {
        log.Printf("Info read error: %v", readerr)
    }

    const prefix = "/"
    r := rest.MakeWebservice(info, args.jobs, prefix)
    http.Handle(prefix, r)

    addr := fmt.Sprintf("%v:%v", args.bind, args.port)
    httperr := http.ListenAndServe(addr, nil)
    if httperr != nil {
        log.Fatal(httperr)
    }
}

func readArgs() params {
    args := params{}
    flag.UintVar(&args.jobs, "jobs", 1, "Number of concurrent render jobs")
    flag.UintVar(&args.port, "port", 98989, "Port for webservice")
    flag.StringVar(&args.bind, "bind", "127.0.0.1", "Interface to bind against")
    flag.Parse()
    return args
}