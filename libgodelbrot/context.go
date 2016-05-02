package libgodelbrot

import (
    "image"
    "fmt"
    "github.com/johnny-morrice/godelbrot/config"
)

type Renderer interface {
    Render() (*image.NRGBA, error)
}

func MakeRenderer(desc *Info) (Renderer, error) {
    // Check that numerics modes are okay, but do not act on them
    switch desc.NumericsStrategy {
    case config.NativeNumericsMode:
    case config.BigFloatNumericsMode:
    default:
        return nil, fmt.Errorf("Invalid NumericsStrategy: %v", desc.NumericsStrategy)
    }

    renderer := Renderer(nil)
    switch desc.RenderStrategy {
    case config.SequenceRenderMode:
        renderer = makeSequenceFacade(desc)
    case config.RegionRenderMode:
        renderer = makeRegionFacade(desc)
    default:
        return nil, fmt.Errorf("Invalid RenderStrategy: %v", desc.RenderStrategy)
    }

    return renderer, nil
}