package libgodelbrot

import (
    "log"
    "image"
    "functorama.com/demo/region"
    "functorama.com/demo/nativeregion"
    "functorama.com/demo/bigregion"
)

type regionProvider struct {
    regionConfig region.RegionConfig
    factory *regionNumericsFactory
}

var _ region.RegionProvider = (*regionProvider)(nil)

func (provider *regionProvider) RegionConfig() region.RegionConfig {
    return provider.regionConfig
}

func (provider *regionProvider) RegionNumericsFactory() region.RegionNumericsFactory {
    return provider.factory
}

type regionFacade struct {
    *regionProvider
    *baseFacade
    *drawFacade
}

var _ region.RenderApplication = (*regionFacade)(nil)
var _ Renderer = (*regionFacade)(nil)

func makeRegionFacade(desc *Info) *regionFacade {
    req := desc.UserRequest
    baseApp := makeBaseFacade(desc)
    facade := &regionFacade{
        baseFacade: baseApp,
        drawFacade: makeDrawFacade(desc),
    }

    provider := &regionProvider{}
    provider.factory = &regionNumericsFactory{desc, baseApp, provider}
    provider.regionConfig = region.RegionConfig{
        Samples: req.RegionSamples,
        CollapseSize: req.RegionCollapse,
    }

    facade.regionProvider = provider
    return facade
}

func (facade *regionFacade) RegionNumericsFactory() region.RegionNumericsFactory {
    return facade.factory
}

func (facade *regionFacade) Render() (*image.NRGBA, error) {
    renderer := region.Make(facade)
    return renderer.Render()
}

type regionNumericsFactory struct {
    desc *Info
    baseApp *baseFacade
    provider *regionProvider
}

func (factory *regionNumericsFactory) Build() region.RegionNumerics {
    switch factory.desc.NumericsStrategy {
    case NativeNumericsMode:
        app := makeNativeRegionFacade(factory.desc, factory.baseApp, factory.provider)
        nativeApp := nativeregion.Make(app)
        return &nativeApp
    case BigFloatNumericsMode:
        app := makeBigRegionFacade(factory.desc, factory.baseApp, factory.provider)
        bigApp := bigregion.Make(app)
        return &bigApp
    default:
        log.Panic("Invalid NumericsStrategy", factory.desc.NumericsStrategy)
        return nil
    }
}