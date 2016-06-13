package restclient

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "io"
    "runtime/debug"
    "strings"
    "time"
    "github.com/johnny-morrice/godelbrot/rest/protocol"
)

type HttpClient interface {
    Get(string) (HttpResponse, error)
    Post(string, string, io.Reader) (HttpResponse, error)
}

type HttpResponse interface {
    GetBody() io.ReadCloser
    GetStatusCode() int
    GetStatus() string
    GetHeader() map[string][]string
    Write(io.Writer) error
}

type Config struct {
    Addr string
    Port uint
    Prefix string
    Ticktime uint
    Debug bool
    Http HttpClient
}

type RenderResult struct {
    Image io.Reader
    Status protocol.RQGetResp
}

// Client provides an interface to restfulbrot over the web, using a user-defined web interface.
// This allows inclusion of this stock code within static archives for other platforms, or
// gopherjs for javascript.
type Client struct {
    config Config
    zoom bool
    tick *time.Ticker
}

func New(config Config) *Client {
    web := &Client{}
    web.config = config
    return web
}

func (web *Client) RenderCycle(url string, renreq *protocol.RenderRequest) (*RenderResult, error) {
    // Continue zoom or start anew?
    rqurl, err := web.cycstartUrl(url, renreq)
    if err != nil {
        return nil, err
    }
    for {
        rqstat, err := web.Getrq(rqurl)
        if err != nil {
            return nil, err
        }
        switch rqstat.State {

        case "done":
            imgurl := web.Url(rqstat.ImageURL)
            
            r, err := web.Getimag(imgurl)
            if err != nil {
                return nil, err
            }

            resp := &RenderResult{}
            resp.Image = r
            resp.Status = *rqstat
            return resp, nil

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

func (web *Client) cycstartUrl(url string, renreq *protocol.RenderRequest) (string, error) {
    if url == "" {
        newresp, err := web.Newrq(web.Url("renderqueue"), renreq)
        if err != nil {
            return "", err
        }
        return web.Url(newresp.RQStatusURL), nil
    } else {
        if renreq.WantZoom {
            newresp, err := web.Rqzoom(url, renreq)
            if err != nil {
                return "" ,err
            }
            return web.Url(newresp.RQStatusURL), nil
        } else {
            return url, nil
        }
    }
}

func (web *Client) Newrq(url string, renreq *protocol.RenderRequest) (*protocol.RQNewResp, error) {
    buff, werr := jsonr(renreq)
    if werr != nil {
        return nil, addstack(werr)
    }
    resp, err := web.post(url, "application/json", buff)
    if err != nil {
        return nil, err
    }
    if resp.GetStatusCode() != 200 {
        return nil, httpError(resp)
    }
    defer resp.GetBody().Close()
    rqi := &protocol.RQNewResp{}
    derr := web.decode(resp.GetBody(), rqi)
    return rqi, addstack(derr)
}

func (web *Client) Rqzoom(rqurl string, renreq *protocol.RenderRequest) (*protocol.RQNewResp, error) {
    cacheresp, err := web.Getrq(rqurl)
    if err != nil {
        return nil, err
    }
    renreq.Req = cacheresp.NextReq

    if web.config.Debug {
        log.Printf("Next request for %v was %v", rqurl, cacheresp.NextReq)
    }

    return web.Newrq(web.Url("renderqueue"), renreq)
}

func (web *Client) Getrq(url string) (*protocol.RQGetResp, error) {
    resp, err := web.get(url)
    if err != nil {
        return nil, err
    }
    defer resp.GetBody().Close()
    rqi := &protocol.RQGetResp{}
    derr := web.decode(resp.GetBody(), rqi)
    return rqi, addstack(derr)
}

func (web *Client) Getimag(url string) (io.Reader, error) {
    resp, err := web.get(url)
    if err != nil {
        return nil, err
    }
    if resp.GetStatusCode() != 200 {
        return nil, httpError(resp)
    }
    defer resp.GetBody().Close()
    buff := &bytes.Buffer{}
    _, cpyerr := io.Copy(buff, resp.GetBody())
    return buff, addstack(cpyerr)
}

func (web *Client) Url(path string) string {
    config := web.config
    var url string
    if web.config.Prefix == "" {
        url = fmt.Sprintf("http://%v:%v/%v",
            config.Addr, config.Port, path)
    } else {
        url = fmt.Sprintf("http://%v:%v/%v/%v",
                config.Addr, config.Port, config.Prefix, path)
    }
    // Help the user for once.
    if !strings.HasSuffix(url, "/") {
        url += "/"
    }
    return url
}

type httpFunc func () (HttpResponse, error)

func (web *Client) get(url string) (r HttpResponse, err error) {
    f := func () (HttpResponse, error) {
        if web.config.Debug {
            log.Printf("GET %v", url)
        }

        return web.config.Http.Get(url)
    }
    return web.request(f)
}

func (web *Client) post(url, ctype string, body io.Reader) (HttpResponse, error) {
    f := func () (HttpResponse, error) {
        if web.config.Debug {
            log.Printf("POST %v", url)
        }

        return web.config.Http.Post(url, ctype, body)
    }
    return web.request(f)
}

func (web *Client) request(f httpFunc) (HttpResponse, error) {
    var r HttpResponse
    var err error
    web.cautiously(func () {
        r, err = f()

        if web.config.Debug {
            web.reportResponse(r, err)
        }
    })
    return r, err
}

func (web *Client) reportResponse(r HttpResponse, err error) {
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    log.Printf("Status: %v", r.GetStatus())
    ctypeHeads := r.GetHeader()["Content-Type"]
    if len(ctypeHeads) != 1 {
        log.Printf("Bad Content-Type header")
    } else {
        log.Printf("Content-Type: %v", ctypeHeads[0])
    }
}

func (web *Client) cautiously(f func()) {
    if web.tick == nil {
        web.tick = time.NewTicker(time.Duration(web.config.Ticktime) * time.Millisecond)
    } else {
        <-web.tick.C
    }
    f()
}

func (web *Client) decode(r io.Reader, any interface{}) error {
    if web.config.Debug {
        buff := &bytes.Buffer{}
        r = io.TeeReader(r, buff)
        derr := decode(r, any)
        log.Printf("Decoded: %v", buff.String())
        return derr
    }
    return decode(r, any)
}

func httpError(resp HttpResponse) error {
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

func addstack(err error) error {
    if err == nil {
        return nil
    } else {
        return fmt.Errorf("%v\n%v", err, string(debug.Stack()))
    }
}