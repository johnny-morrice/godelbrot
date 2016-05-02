package libgodelbrot

import (
    "log"
    "image"
    "github.com/johnny-morrice/godelbrot/config"
    "github.com/johnny-morrice/godelbrot/internal/sequence"
    "github.com/johnny-morrice/godelbrot/internal/nativesequence"
    "github.com/johnny-morrice/godelbrot/internal/bigsequence"
)

type sequenceFacade struct {
    *baseFacade
    *drawFacade
    factory *sequenceNumericsFactory
}

// sequenceFacade implements a couple of interfaces
var _ sequence.RenderApplication = (*sequenceFacade)(nil)
var _ Renderer = (*sequenceFacade)(nil)

func makeSequenceFacade(info *Info) *sequenceFacade {
    baseApp := makeBaseFacade(info)
    facade := &sequenceFacade{
        baseFacade: baseApp,
        drawFacade: makeDrawFacade(info),
    }
    facade.factory = &sequenceNumericsFactory{info, baseApp}
    return facade
}

func (facade *sequenceFacade) SequenceNumericsFactory() sequence.SequenceNumericsFactory {
    return facade.factory
}

func (facade *sequenceFacade) Render() (*image.NRGBA, error) {
    renderer := sequence.Make(facade)
    return renderer.Render()
}

type sequenceNumericsFactory struct {
    desc *Info
    baseApp *baseFacade
}

func (factory *sequenceNumericsFactory) Build() sequence.SequenceNumerics {
    switch factory.desc.NumericsStrategy {
    case config.NativeNumericsMode:
        specialBase := makeNativeBaseFacade(factory.desc, factory.baseApp)
        nativeApp := nativesequence.Make(specialBase)
        return &nativeApp
    case config.BigFloatNumericsMode:
        specialBase := makeBigBaseFacade(factory.desc, factory.baseApp)
        bigApp := bigsequence.Make(specialBase)
        return &bigApp
    default:
        log.Panic("Invalid NumericsStrategy", factory.desc.NumericsStrategy)
        return nil
    }
}