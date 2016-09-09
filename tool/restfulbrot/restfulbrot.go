package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	lib "github.com/johnny-morrice/godelbrot"
	"github.com/johnny-morrice/godelbrot/rest"
)

type params struct {
	jobs    uint
	port    uint
	bind    string
	debug   bool
	origins string
}

func main() {
	args := readArgs()

	if args.debug {
		log.Println("Running in debug mode")
	}

	info, readerr := lib.ReadInfo(os.Stdin)
	if readerr != nil {
		log.Printf("Info read error: %v", readerr)
	}

	const prefix = ""
	h := rest.MakeWebservice(info, args.jobs, prefix)

	// Add CORS to the handler
	if args.origins != "" {
		h = newcorshandler(h, args)
	}

	// Add logging to the handler
	if args.debug {
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
	allowed map[string]bool
	headers []string
	debug   bool
}

func newcorshandler(h http.Handler, args params) http.Handler {
	ch := corshandler{}
	ch.handler = h
	ch.debug = args.debug
	ch.headers = []string{"Content-Type"}

	ch.allowed = map[string]bool{}
	for _, o := range strings.Split(args.origins, ",") {
		ch.allowed[o] = true
		if ch.debug {
			log.Printf("Allowing CORS for '%v'", o)
		}
	}

	return ch
}

func (ch corshandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ch.cors(w, req)
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	} else {
		ch.handler.ServeHTTP(w, req)
	}
}

func (ch corshandler) cors(w http.ResponseWriter, req *http.Request) {
	originheads := req.Header["Origin"]
	var origin string
	if len(originheads) == 1 {
		origin = originheads[0]
	} else if ch.debug {
		log.Printf("Bad Origin header: %v", strings.Join(originheads, ","))
	}

	if ch.allowed[origin] {
		w.Header()["Access-Control-Allow-Origin"] = []string{origin}
		w.Header()["Access-Control-Allow-Headers"] = ch.headers
		if ch.debug {
			log.Printf("CORS okay for '%v'", origin)
		}
	} else if ch.debug {
		log.Printf("CORS denied for '%v'", origin)
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
	flag.StringVar(&args.origins, "origins", "", "Comma separated CORS client origins")
	flag.BoolVar(&args.debug, "debug", false, "Verbose logging")
	flag.Parse()
	return args
}
