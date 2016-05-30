package protocol

import (
    "github.com/johnny-morrice/godelbrot/config"
)

type RenderRequest struct {
    Req config.Request
    Target config.ZoomBounds
    WantZoom bool
}

type RQNewResp struct {
    RQStatusURL string
}

type RQGetResp struct {
    CreateTime int64
    CompleteTime int64
    State string
    Error string
    NextReq config.Request
    ImageURL string
    ThisUrl string
}