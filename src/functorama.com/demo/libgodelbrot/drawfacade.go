package libgodelbrot

import (
    "log"
    "image"
    "functorama.com/demo/draw"
)

type DrawFacade struct {
    picture *image.NRGBA
    colors draw.Palette
}

var _ draw.DrawingContext = (*DrawFacade)(nil)


func CreateDrawFacade(info *RenderInfo) *DrawFacade {
    facade := &DrawFacade{}
    facade.colors = createPalette(info)
    facade.picture = createImage(info)
    return facade
}

func createImage(info *RenderInfo) *image.NRGBA {
    desc := info.Desc
    bounds := image.Rectangle{
        Min: image.ZP,
        Max: image.Point{
            X: desc.ImageWidth,
            Y: desc.ImageHeight,
        },
    }
    return image.NewNRGBA(bounds)
}

func createStoredPalette(info *RenderInfo) draw.Palette {
    palettes := map[string]PaletteFactory{
        "redscale": NewRedscalePalette,
        "pretty":   NewPrettyPalette,
    }
    code := info.UserDescription.PaletteCode
    found := palettes[code]
    if found == nil {
        log.Panic("Unknown palette:", code)
    }
    return found
}

func createPalette(info *RenderInfo) draw.Palette {
    desc := info.UserDescription
    // We are planning more types of palettes soon
    switch desc.PaletteType {
    case StoredPalette:
        return createStoredPalette(desc.PaletteCode)
    default:
        log.Panic("Unknown palette kind:", desc.PaletteType)
    }
    return nil
}
