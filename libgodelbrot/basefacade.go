package libgodelbrot

import (
    "github.com/johnny-morrice/godelbrot/internal/base"
)

type baseFacade struct {
    config base.BaseConfig

    pictureWidth uint
    pictureHeight uint
}

var _ base.RenderApplication = (*baseFacade)(nil)

func makeBaseFacade(desc *Info) *baseFacade {
    req := &desc.UserRequest
    facade := &baseFacade{}
    facade.config = base.BaseConfig{
        IterateLimit: req.IterateLimit,
        DivergeLimit: req.DivergeLimit,
        FixAspect: req.FixAspect,
    }
    facade.pictureWidth = req.ImageWidth
    facade.pictureHeight = req.ImageHeight
    return facade
}

func (base *baseFacade) BaseConfig() base.BaseConfig {
    return base.config
}

func (base *baseFacade) PictureDimensions() (uint, uint) {
    return base.pictureWidth, base.pictureHeight
}