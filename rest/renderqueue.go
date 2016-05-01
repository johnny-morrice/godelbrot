package rest

import (
    "bytes"
    "crypto/md5"
    "encoding/base64"
    "encoding/gob"
    "encoding/json"
    "fmt"
    "io"
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

    type UserPacket struct {
        WantZoom bool
        Target lib.ZoomTarget
        Info lib.UserInfo
    }
    userpkt := UserPacket{
        WantZoom: pkt.wantzoom,
        Target: pkt.target,
        Info: *lib.Friendly(&pkt.info),
    }

    // Use json in debug mode but gob otherwise.
    serialize := func (encfactory func (io.Writer) func(interface{}) error) []byte {
        buff := &bytes.Buffer{}
        encoder := encfactory(buff)
        err := encoder(userpkt)
        if err != nil {
            panic(fmt.Sprintf("Render packet should serialize: %v", err))
        }
        return buff.Bytes()
    }
    var dat []byte
    if __DEBUG {
        dat = serialize(func (w io.Writer) func(interface{}) error {
            enc := json.NewEncoder(w)
            return enc.Encode
        })
    } else {
        dat = serialize(func (w io.Writer) func (interface{}) error {
            enc := gob.NewEncoder(w)
            return enc.Encode
        })
    }
    hsh := md5.Sum(dat)
    buff64 := &bytes.Buffer{}
    enc64 := base64.NewEncoder(base64.URLEncoding, buff64)
    enc64.Write(hsh[:])
    rqi.code = hashcode(buff64.String())

    debugf("rqi with packet representation %v has code %v", string(dat), rqi.code)
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
    var elapsed time.Duration
    var pkt renderpacket
    var err string
    writeM(rqi.mutex, func () {
        rqi.completetime = time.Now()
        elapsed = rqi.completetime.Sub(rqi.createtime)
        state = rqi.state
        pkt = rqi.pkt
        err = rqi.err
    })
    switch state {
    case __DONE:
        log.Printf("rqitem rendered OK after %v", elapsed)
    case __ERROR:
        log.Printf("rqitem error after %v: %v", elapsed, err)
    default:
        panic(fmt.Sprintf("rq completed after %v with bad state (%v): %v",
                elapsed, state, pkt))
    }
}

func (rqi *rqitem) mkbuffs() (*renderbuffers, error) {
    buffs := &renderbuffers{}
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
    debugf("Queing packet %v", code)

    _, present := rq.ca.get(code)
    if present {
        debugf("Deduplicated packet %v", code)
        return code
    }

    debugf("Storing packet %v", code)
    rq.ca.put(rqi)
    go func() {
        rq.sysdraw(rqi, pkt)
    }()

    return code
}

func (rq *renderqueue) sysdraw(rqi *rqitem, pkt *renderpacket) {
    code := hashcode("")
    if __DEBUG {
        code = rqi.hash()
    }
    var zoomArgs []string
    if pkt.wantzoom {
        zoomArgs = process.ZoomArgs(pkt.target)
    }

    buffs, berr := rqi.mkbuffs()

    if berr != nil {
        rqi.fail("failed input buffer")
        log.Printf("Buffer input error: %v, for packet %v (%v)",
            berr, rqi.hash(), rqi.packet())
        return
    }
    renderErr := rq.rs.render(buffs, zoomArgs)
    debugf("Rendered packet %v", code)

    // Copy any stderr messages
    buffs.logReport()

    if renderErr != nil {
        rqi.fail("failed render")
        log.Printf("Render error: %v, for packet %v (%v)",
            renderErr, rqi.hash(), rqi.packet())
        return
    }

    nextinfo := new(lib.Info)
    *nextinfo = rqi.packet().info
    zoominfo, infoerr := lib.ReadInfo(&buffs.nextinfo)
    if infoerr != nil {
        rqi.fail("failed output buffer")
        log.Printf("Buffer output error: %v, for packet %v (%v)",
            infoerr, rqi.hash(), rqi.packet())
        return
    }
    nextinfo = zoominfo

    rq.ic.put(rqi.hash(), buffs.png.Bytes())
    rqi.done(nextinfo)
}
