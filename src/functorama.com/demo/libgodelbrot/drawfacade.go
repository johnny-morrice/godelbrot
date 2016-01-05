package libgodelbrot

import (
    "log"
    "image"
    "functorama.com/demo/draw"
)

type drawFacade struct {
    picture *image.NRGBA
    colors draw.Palette
}

var _ draw.DrawingContext = (*drawFacade)(nil)
var _ draw.ContextProvider = (*drawFacade)(nil)

func (facade *drawFacade) DrawingContext() draw.DrawingContext {
    return facade
}

func (facade *drawFacade) Colors() draw.Palette {
    return facade.colors
}

func (facade *drawFacade) Picture() *image.NRGBA {
    return facade.picture
}

func makeDrawFacade(desc *Info) *drawFacade {
    facade := &drawFacade{}
    facade.colors = createPalette(desc)
    facade.picture = createImage(desc)
    return facade
}

func createImage(desc *Info) *image.NRGBA {
    req := desc.UserRequest
    bounds := image.Rectangle{
        Min: image.ZP,
        Max: image.Point{
            X: int(req.ImageWidth),
            Y: int(req.ImageHeight),
        },
    }
    return image.NewNRGBA(bounds)
}

func createStoredPalette(desc *Info) draw.Palette {
    palettes := map[string]draw.PaletteFactory{
        "redscale": draw.NewRedscalePalette,
        "pretty":   draw.NewPrettyPalette,
    }
    code := desc.UserRequest.PaletteCode
    found := palettes[code]
    if found == nil {
        log.Panic("Unknown palette:", code)
    }
    return found(desc.UserRequest.IterateLimit)
}

func createPalette(desc *Info) draw.Palette {
    req := desc.UserRequest
    // We are planning more types of palettes soon
    switch req.PaletteType {
    case StoredPalette:
        return createStoredPalette(desc)
    default:
        log.Panic("Unknown palette kind:", req.PaletteType)
    }
    return nil
}
