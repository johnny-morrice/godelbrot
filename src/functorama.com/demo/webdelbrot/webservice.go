package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "time"
    "image/png"
    "log"
    "functorama.com/demo/libgodelbrot"
)

type webCommand uint

const (
    render = webCommand(iota)
    displayImage = webCommand(iota)
)

type webRenderParameters struct {
    imageWidth uint
    imageHeight uint
    realMin float64
    realMax float64
    imagMin float64
    imagMax float64
}

type renderMetadata struct {
    renderer string
    palette string
    renderDuration string
}

type godelbrotPacket struct {
    command webCommand
    renderArgs webRenderParameters
    renderMetadata
}

const godelbrotHeader string = "X-Godelbrot-Packet"

func makeWebserviceHandler() func(http.ResponseWriter, *http.Request) {
    // For now, always use concurrent render and pretty palette
    iterateLimit := uint8(255)
    palette := libgodelbrot.NewPrettyPalette(iterateLimit)
    renderer := libgodelbrot.ConcurrentRegionRender
    threads := libgodelbrot.DefaultRenderThreads()

    baseMetadata := renderMetadata{
        renderer: "ConcurrentRegionRender",
        palette: "Pretty",
    }

    return func (w http.ResponseWriter, req *http.Request) {
        metadata := baseMetadata
        jsonPacket := req.Header.Get(godelbrotHeader)

        if len(jsonPacket) == 0 {
            http.Error(w, fmt.Sprintf("No data found in header '%v'", godelbrotHeader), 400)
        }
        userPacket := godelbrotPacket{}
        jsonError := json.Unmarshal([]byte(jsonPacket), &userPacket)

        if jsonError != nil {
            http.Error(w, fmt.Sprintf("Invalid JSON packet: %v", jsonError), 400)
        }

        args := userPacket.renderArgs

        params := libgodelbrot.RenderParameters{
            IterateLimit: iterateLimit,
            DivergeLimit: 4.0,
            Width: args.imageWidth,
            Height: args.imageHeight,
            RegionCollapse: 2,
            Frame: libgodelbrot.CornerFrame,
            TopLeft: complex(args.realMin, args.imagMax),
            BottomRight: complex(args.realMax, args.imagMin),
            RenderThreads: threads,
            BufferSize: libgodelbrot.DefaultBufferSize,
            FixAspect: true,
        }

        config := params.Configure()

        t0 := time.Now()
        pic, renderError := renderer(config, palette)
        t1 := time.Now()
        metadata.renderDuration = t1.Sub(t0).String()

        if renderError != nil {
            log.Fatal(fmt.Sprintf("Render error: %v", renderError))
        }

        responsePacket := godelbrotPacket{
            command: displayImage,
            renderMetadata: metadata,
        }

        responseHeaderPacket, marshalError := json.Marshal(responsePacket)

        if marshalError != nil {
            http.Error(w, fmt.Sprintf("Error marshalling response header: %v", marshalError), 500)
        }

        // Respond to the request
        w.Header().Set("Content-Type", "image/png")
        w.Header().Set(godelbrotHeader, string(responseHeaderPacket))
        pngError := png.Encode(w, pic)

        if pngError != nil {
            log.Fatal("Error encoding PNG: %v", pngError)
        }
    }
}
