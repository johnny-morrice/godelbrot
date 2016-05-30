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
    "github.com/gorilla/mux"
    "github.com/johnny-morrice/godelbrot/config"
    "github.com/johnny-morrice/godelbrot/rest/protocol"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

func validate(renreq *protocol.RenderRequest) error {
    if renreq.Req.ImageWidth < 1 || renreq.Req.ImageHeight < 1 {
        return errors.New("Invalid Req")
    }

    validerr := renreq.Target.Validate();
    if renreq.WantZoom && validerr != nil {
        return errors.New("Invalid Target")
    }

    if !renreq.WantZoom && validerr == nil {
        return errors.New("False WantZoom yet valid Target")
    }


    return nil
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

func MakeWebservice(baseinfo *lib.Info, concurrent uint, prefix string) http.Handler {
    ws := &webservice{}
    ws.baseinfo = *baseinfo
    ws.prefix = prefix

    ws.rq = makeRenderQueue(concurrent)

    r := mux.NewRouter()
    r.HandleFunc("/renderqueue/", ws.enterRQHandler).Methods("POST")
    r.HandleFunc("/renderqueue/{rqcode}/", ws.getRQHandler).Methods("GET")
    r.HandleFunc("/image/{rqcode}/", ws.getImageHandler).Methods("GET")
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

    resp := &protocol.RQGetResp{}

    // This ugly read is candidate for encapsulation in rqitem
    var completetime int64
    var rqerr string
    var state rqstate
    var nextreq config.Request
    var code hashcode
    readM(rqi.mutex, func () {
        resp.CreateTime = rqi.createtime.Unix()
        if rqi.state == __DONE || rqi.state == __ERROR {
            completetime = rqi.completetime.Unix()
        }
        rqerr = rqi.err
        state = rqi.state
        nextreq = rqi.pkt.info.UserRequest
        code = rqi.code
    })

    resp.ThisUrl = ws.prefixed("renderqueue/%v/", code)

    switch state {
    case __DONE:
        resp.State = "done"
        resp.CompleteTime = completetime
        resp.NextReq = nextreq
        resp.ImageURL = ws.prefixed("image/%v/", code)
    case __ERROR:
        resp.State = "error"
        resp.CompleteTime = completetime
        resp.Error = rqerr
    case __WAIT:
        resp.State = "wait"
    default:
        panic(fmt.Sprintf("Unknown state: %v", state))
    }

    return s.serveJson(resp)
}

func (ws *webservice) enterRQ(s session) error {
    dec := json.NewDecoder(s.req.Body)
    renreq := &protocol.RenderRequest{}
    jsonerr := dec.Decode(renreq)

    if jsonerr != nil {
        err := s.httpError(fmt.Sprintf("Invalid JSON"), 400)
        log.Println(err)
        return jsonerr
    }

    validerr := validate(renreq)
    if validerr != nil {
        err := s.httpError(fmt.Sprintf("Invalid render request: %v", validerr), 400)
        log.Println(err)
        return validerr
    }

    // Defaults overwrite where appropriate (security concern)

    info := ws.mergeInfo(renreq)
    target := ws.makeTarget(renreq)

    pkt := &renderpacket{}
    pkt.wantzoom = renreq.WantZoom
    pkt.info = *info
    pkt.target = target

    code := ws.rq.enqueue(pkt)
    resp := &protocol.RQNewResp{}
    resp.RQStatusURL = ws.prefixed("renderqueue/%v/", code)
    
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

func (ws *webservice) prefixed(format string, args... interface{}) string {
    if ws.prefix == "" {
        return fmt.Sprintf(format, args...)
    } else {
        // TODO inefficient
        more := []interface{}{
            ws.prefix,
        }
        args = append(more, args...)
        return fmt.Sprintf("%v/" + format, args)
    }
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
func (ws *webservice) makeTarget(renreq *protocol.RenderRequest) config.ZoomTarget {
    req := ws.baseinfo.UserRequest
    dyn := req.Renderer == config.AutoDetectRenderMode
    dyn = dyn && req.Numerics == config.AutoDetectNumericsMode

    target := config.ZoomTarget{}

    target.ZoomBounds = renreq.Target
    target.UpPrec = dyn
    target.Reconfigure = dyn
    target.Frames = 1
    return target
}

func (ws *webservice) mergeInfo(renreq *protocol.RenderRequest) *lib.Info {
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