package libgodelbrot

import (
    "image"
    "fmt"
)

type Renderer interface {
    Render() (*image.NRGBA, error)
}

func MakeRenderer(desc *Info) (Renderer, error) {
    // Check that numerics modes are okay, but do not act on them
    switch desc.NumericsStrategy {
    case NativeNumericsMode:
    case BigFloatNumericsMode:
    default:
        return nil, fmt.Errorf("Invalid NumericsStrategy: %v", desc.NumericsStrategy)
    }

    renderer := Renderer(nil)
    switch desc.RenderStrategy {
    case SequenceRenderMode:
        renderer = makeSequenceFacade(desc)
    case RegionRenderMode:
        renderer = makeRegionFacade(desc)
    case SharedRegionRenderMode:
        renderer = makeSharedRegionFacade(desc)
    default:
        return nil, fmt.Errorf("Invalid RenderStrategy: %v", desc.RenderStrategy)
    }

    return renderer, nil
}