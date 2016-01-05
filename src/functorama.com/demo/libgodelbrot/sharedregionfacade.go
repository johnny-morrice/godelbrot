package libgodelbrot

import (
    "log"
    "image"
    "functorama.com/demo/region"
    "functorama.com/demo/sharedregion"
    "functorama.com/demo/bigsharedregion"
    "functorama.com/demo/nativesharedregion"
)

type sharedRegionFacade struct {
    *baseFacade
    *drawFacade
    *sharedRegionProvider
}

var _ sharedregion.RenderApplication = (*sharedRegionFacade)(nil)
var _ Renderer = (*sharedRegionFacade)(nil)

func makeSharedRegionFacade(desc *Info) *sharedRegionFacade {
    req := desc.UserRequest
    baseApp := makeBaseFacade(desc)
    regionConfig := region.RegionConfig{
        GlitchSamples: req.GlitchSamples,
        CollapseSize: req.RegionCollapse,
    }
    sharedConfig := sharedregion.SharedRegionConfig{
        Jobs: desc.UserRequest.Jobs,
    }
    provider := &sharedRegionProvider{
        sharedConfig: sharedConfig,
    }
    provider.regionConfig = regionConfig
    facade := &sharedRegionFacade{
        baseFacade: baseApp,
        drawFacade: makeDrawFacade(desc),
        sharedRegionProvider: provider,
    }
    provider.factory = &sharedRegionFactory{desc, baseApp, provider}
    return facade
}

func (facade *sharedRegionFacade) SharedRegionNumericsFactory() sharedregion.SharedRegionFactory {
    return facade.factory
}

func (facade *sharedRegionFacade) Render() (*image.NRGBA, error) {
    renderer := sharedregion.Make(facade)
    return renderer.Render()
}

type sharedRegionProvider struct {
    regionProvider
    sharedConfig sharedregion.SharedRegionConfig
    factory *sharedRegionFactory
}

var _ sharedregion.SharedProvider = (*sharedRegionProvider)(nil)
var _ region.RegionProvider = (*sharedRegionProvider)(nil)

func (provider *sharedRegionProvider) SharedRegionConfig() sharedregion.SharedRegionConfig {
    return provider.sharedConfig
}

func (provider *sharedRegionProvider) SharedRegionFactory() sharedregion.SharedRegionFactory {
    return provider.factory
}

type sharedRegionFactory struct {
    desc *Info
    baseApp *baseFacade
    provider *sharedRegionProvider
}

func (factory *sharedRegionFactory) Build() sharedregion.SharedRegionNumerics {
    switch factory.desc.NumericsStrategy {
    case NativeNumericsMode:
        app := makeNativeSharedRegionFacade(factory.desc, factory.baseApp, factory.provider)
        return nativesharedregion.Make(app)
    case BigFloatNumericsMode:
        app := makeBigSharedRegionFacade(factory.desc, factory.baseApp, factory.provider)
        return bigsharedregion.Make(app)
    default:
        log.Panic("Unknown numerics mode", factory.desc.NumericsStrategy)
        return nil
    }
}