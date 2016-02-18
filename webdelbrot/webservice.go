package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "time"
    "image/png"
    "log"
    "functorama.com/demo/libgodelbrot"
    "bytes"
    "runtime"
)

type WebCommand string
const (
    render = WebCommand("render")
    displayImage = WebCommand("displayImage")
)

type WebRenderParameters struct {
    ImageWidth uint
    ImageHeight uint
    RealMin float64
    RealMax float64
    ImagMin float64
    ImagMax float64
}

type RenderMetadata struct {
    Renderer string
    Palette string
    RenderDuration string
}

type GodelbrotPacket struct {
    Command WebCommand
    Render WebRenderParameters
    Metadata RenderMetadata
}

const godelbrotHeader string = "X-Godelbrot-Packet"
const godelbrotGetParam string = "godelbrotPacket"

type queueCommand uint

const (
    queueRender = queueCommand(iota)
    queueStop = queueCommand(iota)
)

type renderQueueItem struct {
    command queueCommand
    w http.ResponseWriter
    req *http.Request
    complete chan<- bool
}

func launchRenderService() (func(http.ResponseWriter, *http.Request), chan<- renderQueueItem) {
    input := make(chan renderQueueItem, libgodelbrot.Kilo)

    go handleRenderRequests(input)

    return httpChanWriter(input), input
}

func httpChanWriter(input chan<- renderQueueItem) func(http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, req *http.Request) {
        done := make(chan bool)
        input <- renderQueueItem{
            command: queueRender, 
            w: w, 
            req: req,
            complete: done,
        }
        // Block until rendering is complete
        <- done
    }
}

type httpRenderBase struct {
    queueItem renderQueueItem
    renderParams libgodelbrot.RenderParameters
    palette libgodelbrot.Palette
    renderer libgodelbrot.Renderer
    metadata RenderMetadata 
}

func handleRenderRequests(input <-chan renderQueueItem) {
    processed := false
    run := true

    // Values constant across all renders
    iterateLimit := uint8(255)
    renderBase := httpRenderBase{
        renderParams: libgodelbrot.RenderParameters{
            IterateLimit: iterateLimit,
            DivergeLimit: 4.0,
            RegionCollapse: 2,
            Frame: libgodelbrot.CornerFrame,
            RenderThreads: libgodelbrot.DefaultRenderThreads(),
            BufferSize: libgodelbrot.DefaultBufferSize,
            FixAspect: true,
        },
        metadata: RenderMetadata{
            Renderer: "ConcurrentRegionRender",
            Palette: "Pretty",
        },
        palette: libgodelbrot.NewPrettyPalette(iterateLimit),
        renderer: libgodelbrot.ConcurrentRegionRender,
    }

    for run {
        select {
        case queueItem := <- input:
            switch queueItem.command {
            case queueRender:
                renderBase.queueItem = queueItem
                renderBase.render()
                processed = true
            case queueStop:
                run = false
            default:
                panic(fmt.Sprintf("Unknown queueCommand: %v", queueItem.command))
            } 
        default:
            if processed {
                // No renders waiting...
                // A good time to force GC!
                runtime.GC()
                processed = false 
            }
        }
    }
}

func (renderBase httpRenderBase) render() {

    jsonPacket := renderBase.queueItem.req.URL.Query().Get(godelbrotGetParam)

    if len(jsonPacket) == 0 {
        http.Error(renderBase.queueItem.w, fmt.Sprintf("No data found in parameter '%v'", godelbrotHeader), 400)
        return
    }

    jsonBytes := []byte(jsonPacket)
    userPacket := GodelbrotPacket{}
    jsonError := json.Unmarshal(jsonBytes, &userPacket)

    if jsonError != nil {
        http.Error(renderBase.queueItem.w, fmt.Sprintf("Invalid JSON packet: %v", jsonError), 400)
        return
    }

    args := userPacket.Render

    if args.ImageWidth == 0 || args.ImageHeight == 0 {
        http.Error(renderBase.queueItem.w, "ImageHeight and ImageWidth cannot be 0", 422)
        return
    }

    renderParams := renderBase.renderParams
    renderParams.Width = args.ImageWidth
    renderParams.Height = args.ImageHeight
    renderParams.TopLeft = complex(args.RealMin, args.ImagMax)
    renderParams.BottomRight = complex(args.RealMax, args.ImagMin)

    config := renderParams.Configure()

    t0 := time.Now()
    pic, renderError := renderBase.renderer(config, renderBase.palette)
    t1 := time.Now()

    if renderError != nil {
        log.Fatal(fmt.Sprintf("Render error: %v", renderError))
    }

    buff := bytes.Buffer{}
    pngError := png.Encode(&buff, pic)

    if pngError != nil {
        log.Println("Error encoding PNG: ", pngError)
        http.Error(renderBase.queueItem.w, fmt.Sprintf("Error encoding PNG: %v", pngError), 500)
    }

    // Craft response
    responsePacket := GodelbrotPacket{
        Command: displayImage,
        Metadata: renderBase.metadata,
    }
    responsePacket.Metadata.RenderDuration = t1.Sub(t0).String()

    log.Println("Render complete in ", responsePacket.Metadata.RenderDuration, 
        "plane co-ords: ", config.TopLeft, config.BottomRight)

    responseHeaderPacket, marshalError := json.Marshal(responsePacket)

    if marshalError != nil {
        http.Error(renderBase.queueItem.w, fmt.Sprintf("Error marshalling response header: %v", marshalError), 500)
    }

    // Respond to the request
    renderBase.queueItem.w.Header().Set("Content-Type", "image/png")
    renderBase.queueItem.w.Header().Set(godelbrotHeader, string(responseHeaderPacket))

    // Write image buffer as http response
    renderBase.queueItem.w.Write(buff.Bytes())

    // Notify that rendering is complete
    renderBase.queueItem.complete <- true
}
