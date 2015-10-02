package main

import (
    "fmt"
    "http"
    "encoding/json"
    "time"
    "bytes"
    "image/png"
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
    imageMin float64
    imageMax float64
}

type renderMetadata {
    renderer string
    palette string
    duration string
}

type godelbrotPacket struct {
    command webCommand,
    renderArgs webRenderParameters
    image []byte
    renderMetadata renderMetadata
}

func makeWebserviceHandler() {
    // For now, always use concurrent render and pretty palette
    iterateLimit := 255
    palette := libgodelbrot.NewPrettyPalette(iterateLimit)
    renderer := libgodelbrot.ConcurrentRegionRender
    threads := libgodelbrot.DefaultRenderThreads()

    baseMetadata := renderMetadata{
        renderer: "ConcurrentRegionRender",
        palette: "Pretty",
    }

    return func (w http.ResponseWriter, req *http.Request) {
        metadata := baseMetadata
        jsonPacket := req.FormValue("godelbrotPacket")

        if len(jsonPacket) == 0 {
            http.Error(w, "No data in form value 'godelbrotPacket'", 400)
        }
        userPacket := godelbrotPacket{}
        jsonError := json.Unmarshal([]byte(jsonPacket), &userPacket)

        if jsonError != nil {
            http.Error(w, fmt.Sprintf("Invalid JSON packet: %v", jsonError), 400)
        }

        params := libgodelbrot.RenderParameters{
            IterateLimit: iterateLimit,
            DivergeLimit: 4.0,
            Width: userPacket.imageWidth,
            Height: userPacket.imageHeight,
            RegionCollapse: 2,
            Frame: libgodelbrot.CornerFrame,
            TopLeft: complex(userPacket.realMin, userPacket.imagMax),
            BottomRight: complex(userPacket.realMax, userPacket.imageMin,
            RenderThreads: threads,
            BufferSize: libgodelbrot.DefaultBufferSize,
            FixAspect: true
        }

        config := params.Configure()

        t0 := time.Now()
        pic, renderError := renderer(config, palette)
        t1 := time.Now()
        metadata.renderDuration := t1.Sub(t0).String()

        picBuffer := bytes.Buffer{}
        pngError := png.Encode(picBuffer, pic)

        if pngError != nil {
            http.Error(w, fmt.Sprintf("Error encoding PNG: %v", pngError), 500)
        }

        responsePacket := godelbrotPacket{
            command: displayImage,
            image: picBuffer.Bytes(),
            renderMetadata: metadata,
        }
    }
}
