package restclient

import (
    "fmt"
    "io"
    "time"
)

type HttpRequest struct {
    Method string
    Url string
    Ctype string // TODO use struct tag
    Body io.Reader
}

type WebResult struct {
    Response HttpResponse
    Error error
}

type RestIO struct {
    Request HttpRequest
    Response chan<- WebResult
}

// UserControl is loosely coupled request/response mechanism
type UserControl struct {
    timeout uint
    dockets chan RestIO
}

var _ HttpClient = (*UserControl)(nil)

func NewUserControl(timeout uint) *UserControl {
    uc := &UserControl{}

    uc.timeout = timeout
    uc.dockets = make(chan RestIO)

    return uc
}

func (uc *UserControl) Get(url string) (HttpResponse, error) {
    req := HttpRequest{}

    req.Method = "GET"
    req.Url = url

    return uc.send(req)
}

func (uc *UserControl) Post(url string, ctype string, body io.Reader) (HttpResponse, error) {
    req := HttpRequest{}

    req.Method = "POST"
    req.Url = url
    req.Ctype = ctype
    req.Body = body

    return uc.send(req)
}

// Communicate with the restclient library asynchronously using Comms()
func (uc *UserControl) Comms() <-chan RestIO {
    return uc.dockets
}

func (uc *UserControl) send(req HttpRequest) (HttpResponse, error) {
    docket := RestIO{}

    docket.Request = req
    resch := make(chan WebResult)
    docket.Response = resch

    go func() {
        uc.dockets<- docket
    }()

    var result WebResult
    
    if uc.timeout > 0 {
        r, err := uc.recvtimeout(resch)
        if err != nil {
            return nil, err
        }
        result = r
    } else {
        result = <-resch
    }

    return result.Response, result.Error

}

func (uc *UserControl) recvtimeout(resch <-chan WebResult) (WebResult, error) {
    timeout := make(chan bool, 1)
    after := time.Duration(uc.timeout)

    go func() {
        time.Sleep(after)
        timeout<- true
    }()

    select {
    case res := <-resch:
        return res, nil
    case <-timeout:
        return WebResult{}, fmt.Errorf("Timeout after %v", after)
    }
}
