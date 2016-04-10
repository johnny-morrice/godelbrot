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
    packet renderpacket
    code hashcode
    state rqstate
    err string
    nextinfo lib.Info
    mutex sync.RWMutex
}

func makeRqitem(pkt *renderpacket) *rqitem {
    rqi := &rqitem{}
    rqi.packet = *pkt
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

func (rqi *rqitem) hash() hashcode {
    rqi.mutex.RLock()
    code := rqi.code
    rqi.mutex.RUnlock()
    return code
}

func (rqi *rqitem) done(nextinfo *lib.Info) {
    rqi.nextinfo = *nextinfo
    rqi.state = __DONE
    rqi.logcomplete()
}

func (rqi *rqitem) fail(msg string) {
    rqi.err = msg
    rqi.logcomplete()
}

func (rqi *rqitem) logcomplete() {
    rqi.completetime = time.Now()
    elapsed := rqi.completetime.Sub(rqi.createtime)
    milli := elapsed * time.Millisecond
    switch rqi.state {
    case __DONE:
        log.Printf("rqitem rendered OK after %v milliseconds", milli)
    case __ERROR:
        log.Printf("rqitem error after %v milliseconds: %v", milli)
    default:
        panic(fmt.Sprintf("rq completed after %v milliseconds with bad state (%v): %v",
                milli, rqi.state, rqi.packet))
    }
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

    _, present := rq.ca.get(rqi.hash())
    if present {
        return rqi.hash()
    }

    rq.ca.put(rqi)
    go func() {
        rq.sysdraw(rqi, pkt)
    }()

    return rqi.hash()
}

func (rq *renderqueue) sysdraw(rqi *rqitem, pkt *renderpacket) {
    var zoomArgs []string
    if pkt.wantzoom {
        zoomArgs = process.ZoomArgs(pkt.target)
    }

    rqi.mutex.RLock()
    info := rqi.packet.info
    rqi.mutex.Lock()
    buffs := renderBuffers{}
    bufferr := buffs.input(&info)
    if bufferr != nil {
        rqi.fail("failed input buffer")
        log.Printf("Buffer input error: %v, for packet %v", bufferr, rqi.packet)
        return
    }
    renderErr := rq.rs.render(buffs, zoomArgs)

    // Copy any stderr messages
    buffs.logReport()

    if renderErr != nil {
        rqi.fail("failed render")
        log.Printf("Render error: %v, for packet %v", renderErr, rqi.packet)
        return
    }

    nextinfo, infoerr := lib.ReadInfo(&buffs.nextinfo)
    if infoerr != nil {
        rqi.fail("failed output buffer")
        log.Printf("Buffer output error: %v, for packet %v", infoerr, rqi.packet)
        return
    }

    rq.ic.put(rqi.hash(), buffs.png.Bytes())
    rqi.mutex.Lock()
    defer rqi.mutex.Unlock()
    rqi.done(nextinfo)
}
