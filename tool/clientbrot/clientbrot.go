package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	lib "github.com/johnny-morrice/godelbrot"
	"github.com/johnny-morrice/godelbrot/config"
	"github.com/johnny-morrice/godelbrot/rest/protocol"
	rc "github.com/johnny-morrice/godelbrot/restclient"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	args := readArgs()

	if (args.cycle && args.getrq == "") || args.newrq {
		if args.config.Debug {
			log.Println("Reading stdin...")
		}
		info, ierr := lib.ReadInfo(os.Stdin)
		fatalguard(ierr)
		args.req = info.GenRequest()
	}

	shcl := newShellClient(args)

	zoomparts := map[string]bool{
		"xmin": true,
		"xmax": true,
		"ymin": true,
		"ymax": true,
	}

	first := true
	flag.Visit(func(fl *flag.Flag) {
		_, ok := zoomparts[fl.Name]
		if ok {
			shcl.zoom = true
			if args.config.Debug && first {
				log.Printf("Found parameter %v: zooming", fl.Name)
				first = false
			}
		}
	})

	// Ugly
	var r io.Reader
	if args.cycle {
		result, err := shcl.cycle()
		fatalguard(err)
		if args.config.Debug {
			log.Printf("RequestUrl: %v", result.Status.ThisUrl)
		}
		r = result.Image
	} else if args.newrq {
		rqi, err := shcl.newrq()
		fatalguard(err)
		reqr, jerr := jsonr(rqi)
		fatalguard(jerr)
		r = reqr
	} else if args.getrq != "" {
		rqi, err := shcl.getrq()
		fatalguard(err)
		reqr, jerr := jsonr(rqi)
		fatalguard(jerr)
		r = reqr
	} else if args.getimag != "" {
		png, err := shcl.getimag()
		fatalguard(err)
		r = png
	}

	_, cpyerr := io.Copy(os.Stdout, r)
	fatalguard(cpyerr)
}

type shellClient struct {
	rc   *rc.Client
	zoom bool
	args params
}

func newShellClient(args params) *shellClient {
	shcl := &shellClient{}
	shcl.args = args
	hcl := &http.Client{}
	hcl.Timeout = time.Millisecond * time.Duration(shcl.args.timeout)
	args.config.Http = (*goHttp)(hcl)
	shcl.rc = rc.New(args.config)
	return shcl
}

func (shcl *shellClient) cycle() (*rc.RenderResult, error) {
	url := ""
	if shcl.args.getrq != "" {
		url = shcl.rc.Url(fmt.Sprintf("renderqueue/%v", shcl.args.getrq))
	}
	return shcl.rc.RenderCycle(url, shcl.renreq())
}

func (shcl *shellClient) getimag() (io.Reader, error) {
	url := fmt.Sprintf("image/%v", shcl.args.getimag)
	return shcl.rc.Getimag(shcl.rc.Url(url))
}

func (shcl *shellClient) newrq() (*protocol.RQNewResp, error) {
	return shcl.rc.Newrq(shcl.rc.Url("renderqueue"), shcl.renreq())
}

func (shcl *shellClient) renreq() *protocol.RenderRequest {
	renreq := &protocol.RenderRequest{}
	renreq.Req = shcl.args.req
	if shcl.zoom {
		renreq.WantZoom = shcl.zoom
		renreq.Target = shcl.args.zoombox
	}
	return renreq
}

func (shcl *shellClient) getrq() (*protocol.RQGetResp, error) {
	url := shcl.rqurl()
	return shcl.rc.Getrq(url)
}

func (shcl *shellClient) rqurl() string {
	return shcl.rc.Url(fmt.Sprintf("renderqueue/%v", shcl.args.getrq))
}

func fatalguard(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

func readArgs() params {
	args := params{}
	flag.StringVar(&args.config.Addr, "remote", "localhost", "Remote address of restfulbrot service")
	flag.UintVar(&args.config.Port, "port", 9898, "Port of remote service")
	flag.StringVar(&args.config.Prefix, "prefix", "", "Prefix of service URL")
	flag.UintVar(&args.config.Ticktime, "ticktime", 100, "Max one request per tick (milliseconds)")
	flag.BoolVar(&args.config.Debug, "debug", false, "Verbose debug mode")
	flag.UintVar(&args.timeout, "timeout", 1000, "Web request abort timeout (milliseconds)")
	flag.BoolVar(&args.newrq, "newrq", false, "Add new item to render queue (info from stdin)")
	flag.StringVar(&args.getrq, "getrq", "", "Get status of render queue item")
	flag.StringVar(&args.getimag, "getimag", "", "Download fractal render (png to stdout)")
	flag.BoolVar(&args.cycle, "cycle", false,
		"Wait for fractal to render (info from stdin, png to stdout")

	flag.UintVar(&args.zoombox.Xmin, "xmin", 0, "xmin pixel bound")
	flag.UintVar(&args.zoombox.Xmax, "xmax", 0, "xmax pixel bound")
	flag.UintVar(&args.zoombox.Ymin, "ymin", 0, "ymin pixel bound")
	flag.UintVar(&args.zoombox.Ymax, "ymax", 0, "ymax pixel bound")
	flag.Parse()

	// Cycle is default on only if no other operations provided.
	operation := map[string]bool{
		"getrq":   true,
		"newrq":   true,
		"getimag": true,
	}
	notcycle := false
	defcycle := false
	flag.Visit(func(fl *flag.Flag) {
		if fl.Name == "cycle" {
			defcycle = true
		} else {
			notcycle = notcycle || operation[fl.Name]
		}
	})
	args.cycle = (defcycle && args.cycle) || !notcycle

	return args
}

type params struct {
	config  rc.Config
	newrq   bool
	getrq   string
	getimag string
	cycle   bool
	timeout uint

	req config.Request

	zoombox config.ZoomBounds
}

func jsonr(any interface{}) (io.Reader, error) {
	buff := &bytes.Buffer{}
	enc := json.NewEncoder(buff)
	err := enc.Encode(any)
	return buff, err
}
