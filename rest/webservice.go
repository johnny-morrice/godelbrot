package rest

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "strconv"
    "strings"
    "github.com/gorilla/mux"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

const formkey = "godelbrot-packet"

type RenderRequest struct {
    Req lib.Request
    Target lib.ZoomTarget
    WantZoom bool
}

func (rr *RenderRequest) validate() error {
    if rr.Req.ImageWidth < 1 || rr.Req.ImageHeight < 1 {
        return errors.New("Invalid Req")
    }

    validerr := rr.Target.Validate();
    if rr.WantZoom && validerr != nil {
        return errors.New("Invalid Target")
    }

    if !rr.WantZoom && validerr == nil {
        return errors.New("False WantZoom yet valid Target")
    }


    return nil
}

type RQNewResp struct {
    RQStatusUrl string
}

type RQStatusResp struct {
    CreateTime int64
    CompleteTime int64
    State string
    Error string
    NextReq lib.Request
    ImageURL string
}

type session struct {
    w http.ResponseWriter
    req *http.Request
}

func (s session) getMuxVar(field string) string {
    mvars := mux.Vars(s.req)
    return mvars[field]
}

func (s session) httpError(msg string, code int) error {
    http.Error(s.w, msg, code)
    return errors.New(fmt.Sprintf("(%v) %v", s.req.RemoteAddr, msg))
}

func (s session) internalError() error {
    return s.httpError("Internal error", 500)
}

func (s session) serveJson(any interface{}) error {
    s.w.Header().Set("Content-Type", "application/json")
    enc := json.NewEncoder(s.w)
    jsonErr := enc.Encode(any)
    if jsonErr != nil {
        err := s.internalError()
        log.Println(err)
        return jsonErr
    }
    return nil
}

type webservice struct {
    baseinfo lib.Info
    rq renderqueue
    prefix string
}

func MakeWebservice(baseinfo *lib.Info, concurrent uint, prefix string) *mux.Router {
    ws := &webservice{}
    ws.baseinfo = *baseinfo
    ws.prefix = prefix

    ws.rq = makeRenderQueue(concurrent)

    r := mux.NewRouter()
    r.HandleFunc("/renderqueue", ws.enterRQHandler)
    r.HandleFunc("/renderqueue/{rqcode}/", ws.getRQHandler)
    r.HandleFunc("/image/{rqcode}/", ws.getImageHandler)
    return r
}

func (ws *webservice) getImageHandler(w http.ResponseWriter, req *http.Request) {
    withSession(w, req, ws.getImage)
}

func (ws *webservice) enterRQHandler(w http.ResponseWriter, req *http.Request) {
    withSession(w, req, ws.enterRQ)
}

func (ws *webservice) getRQHandler(w http.ResponseWriter, req *http.Request) {
    withSession(w, req, ws.getRQ)
}

func (ws *webservice) getRQ(s session) error {
    input := s.getMuxVar("rqcode")
    rqcode := hashcode(input)
    any, ok := ws.rq.ca.get(rqcode)

    if !ok {
        err := s.httpError(fmt.Sprintf("Invalid code: %v", rqcode), 400)
        return err
    }

    rqi, castok := any.(*rqitem)
    if !castok {
        panic(fmt.Sprintf("Expected type rqitem but received: %v", any))
    }

    resp := &RQStatusResp{}
    resp.CreateTime = rqi.createtime.Unix()

    switch rqi.state {
    case DONE:
        resp.State = "done"
        resp.CompleteTime = rqi.completetime.Unix()
        resp.NextReq = rqi.packet.info.UserRequest
        resp.ImageURL = fmt.Sprintf("%v/image/%v/", ws.prefix, rqi.hash())
    case ERROR:
        resp.State = "error"
        resp.CompleteTime = rqi.completetime.Unix()
        resp.Error = rqi.err
    case WAIT:
        resp.State = "wait"
    default:
        panic(fmt.Sprintf("Unknown state: %v", rqi.state))
    }

    return s.serveJson(resp)
}

func (ws *webservice) enterRQ(s session) error {
    jsonPacket := s.req.FormValue(formkey)

    if len(jsonPacket) == 0 {
        err := s.httpError(fmt.Sprintf("No data found in parameter '%v'", formkey), 400)
        return err
    }

    dec := json.NewDecoder(strings.NewReader(jsonPacket))
    renreq := &RenderRequest{}
    jsonerr := dec.Decode(renreq)

    if jsonerr != nil {
        err := s.httpError(fmt.Sprintf("Invalid JSON"), 400)
        log.Println(err)
        return jsonerr
    }

    validerr := renreq.validate()
    if validerr != nil {
        err := s.httpError(fmt.Sprintf("Invalid render request: %v", validerr), 400)
        log.Println(err)
        return validerr
    }

    // Defaults overwrite where appropriate (security concern)
    ws.safeTarget(renreq)
    info := ws.mergeInfo(renreq)

    pkt := &renderpacket{}
    pkt.wantzoom = renreq.WantZoom
    pkt.info = *info
    pkt.target = renreq.Target

    code := ws.rq.enqueue(pkt)
    resp := &RQNewResp{}
    resp.RQStatusUrl = fmt.Sprintf("%v/renderqueue/%v/", ws.prefix, code)
    return s.serveJson(resp)
}


func (ws *webservice) getImage(s session) error {
    input := s.getMuxVar("rqcode")
    rqcode := hashcode(input)

    png, ok := ws.rq.ic.get(rqcode)

    if !ok {
        err := s.httpError(fmt.Sprintf("Invalid Code: %v", rqcode), 400)
        return err
    }

    buff := bytes.NewBuffer(png)
    // Write image buffer as http response
    s.w.Header().Set("Content-Type", "image/png")
    _, cpyerr := io.Copy(s.w, buff)
    if cpyerr != nil {
        err := s.internalError()
        log.Println(err)
        return cpyerr
    }

    return nil
}

func withSession(w http.ResponseWriter, req *http.Request, handler func (session) error) {
    sess := session{}
    sess.w = w
    sess.req = req
    err := handler(sess)
    if err != nil {
        log.Println(err)
    }
}

// Only allow zoom reconfiguration if autodetection is enabled throughout the base info.
func (ws *webservice) safeTarget(renreq *RenderRequest) {
    req := ws.baseinfo.UserRequest
    dyn := req.Renderer == lib.AutoDetectRenderMode
    dyn = dyn && req.Numerics == lib.AutoDetectNumericsMode

    renreq.Target.UpPrec = dyn
    renreq.Target.Reconfigure = dyn
    renreq.Target.Frames = 1
}

func (ws *webservice) mergeInfo(renreq *RenderRequest) *lib.Info {
    req := ws.baseinfo.UserRequest

    req.ImageWidth = renreq.Req.ImageWidth
    req.ImageHeight = renreq.Req.ImageHeight

    inf := new(lib.Info)
    *inf = ws.baseinfo
    inf.UserRequest = req

    return inf
}

func format(f float64) string {
    return strconv.FormatFloat(f, 'e', -1, 64)
}