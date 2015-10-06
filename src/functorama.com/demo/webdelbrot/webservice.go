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

func makeWebserviceHandler() func(http.ResponseWriter, *http.Request) {
    // For now, always use concurrent render and pretty palette
    iterateLimit := uint8(255)
    palette := libgodelbrot.NewPrettyPalette(iterateLimit)
    renderer := libgodelbrot.ConcurrentRegionRender
    threads := libgodelbrot.DefaultRenderThreads()

    baseMetadata := RenderMetadata{
        Renderer: "ConcurrentRegionRender",
        Palette: "Pretty",
    }

    return func (w http.ResponseWriter, req *http.Request) {
        metadata := baseMetadata
        jsonPacket := req.URL.Query().Get(godelbrotGetParam)

        if len(jsonPacket) == 0 {
            http.Error(w, fmt.Sprintf("No data found in parameter '%v'", godelbrotHeader), 400)
            return
        }

        jsonBytes := []byte(jsonPacket)
        userPacket := GodelbrotPacket{}
        jsonError := json.Unmarshal(jsonBytes, &userPacket)

        if jsonError != nil {
            http.Error(w, fmt.Sprintf("Invalid JSON packet: %v", jsonError), 400)
            return
        }

        args := userPacket.Render

        if args.ImageWidth == 0 || args.ImageHeight == 0 {
            http.Error(w, "ImageHeight and ImageWidth cannot be 0", 422)
            return
        }

        params := libgodelbrot.RenderParameters{
            IterateLimit: iterateLimit,
            DivergeLimit: 4.0,
            Width: args.ImageWidth,
            Height: args.ImageHeight,
            RegionCollapse: 2,
            Frame: libgodelbrot.CornerFrame,
            TopLeft: complex(args.RealMin, args.ImagMax),
            BottomRight: complex(args.RealMax, args.ImagMin),
            RenderThreads: threads,
            BufferSize: libgodelbrot.DefaultBufferSize,
            FixAspect: true,
        }

        config := params.Configure()

        t0 := time.Now()
        pic, renderError := renderer(config, palette)
        t1 := time.Now()
        metadata.RenderDuration = t1.Sub(t0).String()

        if renderError != nil {
            log.Fatal(fmt.Sprintf("Render error: %v", renderError))
        }

        buff := bytes.Buffer{}
        pngError := png.Encode(&buff, pic)

        if pngError != nil {
            log.Println("Error encoding PNG: ", pngError)
            http.Error(w, fmt.Sprintf("Error encoding PNG: %v", pngError), 500)
        }

        responsePacket := GodelbrotPacket{
            Command: displayImage,
            Metadata: metadata,
        }

        log.Println("Render complete in ", metadata.RenderDuration, 
            "plane co-ords: ", config.TopLeft, config.BottomRight)

        responseHeaderPacket, marshalError := json.Marshal(responsePacket)

        if marshalError != nil {
            http.Error(w, fmt.Sprintf("Error marshalling response header: %v", marshalError), 500)
        }

        // Respond to the request
        w.Header().Set("Content-Type", "image/png")
        w.Header().Set(godelbrotHeader, string(responseHeaderPacket))

        // Write image buffer as http response
        w.Write(buff.Bytes())
    }
}
