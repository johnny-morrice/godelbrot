package libgodelbrot

import (
    "functorama.com/demo/base"
)

type BaseFacade struct {
    config base.BaseConfig

    pictureWidth uint
    pictureHeight uint
}

func NewBaseFacade(info *RenderInfo) *BaseFacade {
    desc := &info.UserDescription
    return BaseFacade{
        config: base.BaseConfig{
            IterateLimit: desc.IterateLimit,
            DivergeLimit: desc.DivergeLimit,
            FixAspect: desc.FixAspect,
        },
    }
}

func (base *BaseFacade) BaseConfig() base.BaseConfig {
    return base.config
}

func (base *BaseFacade) PictureDimensions() (uint, uint) {
    return pictureWidth, pictureHeight
}