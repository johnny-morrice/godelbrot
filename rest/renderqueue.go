package rest

import (
    "bytes"
    "crypto/md5"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    "time"
    "sync"
    "github.com/johnny-morrice/godelbrot/process"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

type rqstate uint8

const (
    __WAIT = rqstate(iota)
    __DONE
    __ERROR
)

type rqitem struct {
    createtime time.Time
    completetime time.Time
    pkt renderpacket
    code hashcode
    state rqstate
    err string
    nextinfo lib.Info
    mutex sync.RWMutex
}

func makeRqitem(pkt *renderpacket) *rqitem {
    rqi := &rqitem{}
    rqi.pkt = *pkt
    rqi.createtime = time.Now()

    buff := &bytes.Buffer{}
    enc := json.NewEncoder(buff)
    err := enc.Encode(pkt)
    if err != nil {
        panic("renderpacket should serialize")
    }
    hsh := md5.Sum(buff.Bytes())
    buff64 := &bytes.Buffer{}
    enc64 := base64.NewEncoder(base64.URLEncoding, buff64)
    enc64.Write(hsh[:])
    rqi.code = hashcode(buff64.String())
    return rqi
}

func (rqi *rqitem) packet() renderpacket {
    rqi.mutex.RLock()
    defer rqi.mutex.RUnlock()
    return rqi.pkt
}

func (rqi *rqitem) hash() hashcode {
    rqi.mutex.RLock()
    defer rqi.mutex.RUnlock()
    return rqi.code
}

func (rqi *rqitem) done(nextinfo *lib.Info) {
    writeM(rqi.mutex, func () {
        rqi.nextinfo = *nextinfo
        rqi.state = __DONE
    })
    rqi.logcomplete()
}

func (rqi *rqitem) fail(msg string) {
    writeM(rqi.mutex, func () {
        rqi.err = msg
        rqi.state = __ERROR
    })
    rqi.logcomplete()
}

func (rqi *rqitem) logcomplete() {
    var state rqstate
    var milli time.Duration
    var pkt renderpacket
    var err string
    writeM(rqi.mutex, func () {
        rqi.completetime = time.Now()
        elapsed := rqi.completetime.Sub(rqi.createtime)
        milli = elapsed * time.Millisecond
        state = rqi.state
        pkt = rqi.pkt
        err = rqi.err
    })
    switch state {
    case __DONE:
        log.Printf("rqitem rendered OK after %v milliseconds", milli)
    case __ERROR:
        log.Printf("rqitem error after %v milliseconds: %v", milli, err)
    default:
        panic(fmt.Sprintf("rq completed after %v milliseconds with bad state (%v): %v",
                milli, state, pkt))
    }
}

func (rqi *rqitem) mkbuffs() (renderbuffers, error) {
    buffs := renderbuffers{}
    var bufferr error
    readM(rqi.mutex, func () {
        info := rqi.pkt.info
        bufferr = buffs.input(&info)
    })
    return buffs, bufferr
}

type renderpacket struct {
    wantzoom bool
    target lib.ZoomTarget
    info lib.Info
}

type renderqueue struct {
    ca cache
    rs renderservice
    ic imagecache
}

func makeRenderQueue(concurrent uint) renderqueue {
    rq := renderqueue{}
    rq.ca = makeCache(&rqitem{})
    rq.ic = makeImageCache()
    rq.rs = makeRenderservice(concurrent)
    return rq
}

func (rq *renderqueue) enqueue(pkt *renderpacket) hashcode {
    rqi := makeRqitem(pkt)
    code := rqi.code

    _, present := rq.ca.get(code)
    if present {
        return code
    }

    rq.ca.put(rqi)
    go func() {
        rq.sysdraw(rqi, pkt)
    }()

    return code
}

func (rq *renderqueue) sysdraw(rqi *rqitem, pkt *renderpacket) {
    var zoomArgs []string
    if pkt.wantzoom {
        zoomArgs = process.ZoomArgs(pkt.target)
    }

    buffs, berr := rqi.mkbuffs()

    if berr != nil {
        rqi.fail("failed input buffer")
        log.Printf("Buffer input error: %v, for packet %v", berr, rqi.packet())
        return
    }
    renderErr := rq.rs.render(buffs, zoomArgs)

    // Copy any stderr messages
    buffs.logReport()

    if renderErr != nil {
        rqi.fail("failed render")
        log.Printf("Render error: %v, for packet %v", renderErr, rqi.packet())
        return
    }

    nextinfo, infoerr := lib.ReadInfo(&buffs.nextinfo)
    if infoerr != nil {
        rqi.fail("failed output buffer")
        log.Printf("Buffer output error: %v, for packet %v", infoerr, rqi.packet())
        return
    }

    rq.ic.put(rqi.hash(), buffs.png.Bytes())
    rqi.done(nextinfo)
}
