package libgodelbrot

import (
    "functorama.com/demo/region"
)

type RegionFacade struct {
    *BaseFacade
    *DrawFacade
    regionConfig region.RegionConfig
    factory GodelbrotRegionNumericsFactory
}

var _ region.RenderApplication = (*RegionFacade)(nil)

func NewRegionFacade(info *RenderInfo) *RegionFacade {
    baseApp := NewBaseFacade()
    facade := &RegionFacade{
        BaseFacade: baseApp,
        DrawFacade: NewDrawFacade(info),
    }
    facade.factory = &GodelbrotRegionNumericsFactory{info, baseApp}
    return facade
}

func (facade *RegionFacade) RegionNumericsFactory() region.RegionNumericsFactory {
    return facade.factory
}

type GodelbrotRegionNumericsFactory struct {
    info RenderInfo
    baseApp base.RenderApplication
}

func (factory *GodelbrotRegionNumericsFactory) Build() region.RegionNumerics {
    desc := factory.info.UserDescription
    config := region.RegionConfig{
        GlitchSamples: desc.GlitchSamples,
        Collapse: desc.RegionCollapse,
    }
    switch factory.info.DetectedNumericsMode {
    case NativeNumericsMode:
        nativeBaseApp := CreateNativeBaseFacade(factory.info, factory.baseApp)
        specialized := nativebase.CreateNativeBaseNumerics(nativeBaseApp)
        sequence := nativesequence.NewNativeSequenceNumerics(specialized)
        return nativeregion.NewRegionNumerics(specialized, config, sequence)
    case BigFloatNumericsMode:
        bigBaseApp := CreateBigBaseFacade(factory.info, factory.baseApp)
        specialized := bigbase.CreateBigBaseNumerics(bigBaseApp)
        sequence := bigsequence.NewBigSequenceNumerics(specialized)
        return bigregion.NewRegionNumerics(specialized, config, sequence)
    default:
        log.Panic("Unknown numerics mode", factory.info.DetectedNumericsMode)
    }
}