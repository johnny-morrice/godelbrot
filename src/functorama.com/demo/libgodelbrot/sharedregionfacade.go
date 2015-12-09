package libgodelbrot

type SharedRegionFacade struct {
    *BaseFacade
    *DrawFacade
    regionConfig RegionConfig
    sharedConfig sharedregion.SharedRegionConfig
    factory GodelbrotRegionNumericsFactory
}

var _ sharedregion.RenderApplication = (*SharedRegionFacade)(nil)

func NewSharedRegionFacade(info *RenderInfo) *SharedRegionFacade {
    baseApp := NewBaseFacade()
    regionConfig := region.RegionConfig{
        GlitchSamples: desc.GlitchSamples,
        Collapse: desc.RegionCollapse,
    }
    sharedConfig := sharedregion.SharedRegionConfig{
        BufferSize: desc.ThreadBufferSize,
        Jobs: desc.uint32(Jobs),
    }
    facade := &SharedRegionFacade{
        BaseFacade: baseApp,
        DrawFacade: NewDrawFacade(info),
        sharedConfig: sharedConfig,
        regionConfig: regionConfig,
    }
    facade.factory = &GodelbrotRegionNumericsFactory{info, baseApp, regionConfig, sharedConfig}
    return facade
}

func (facade *SharedRegionFacade) RegionNumericsFactory() sharedregion.SharedRegionNumericsFactory {
    return facade.factory
}

type GodelbrotRegionNumericsFactory struct {
    info RenderInfo
    baseApp base.RenderApplication
    regionConfig RegionConfig
    sharedConfig sharedregion.SharedRegionConfig
}

func (factory *GodelbrotRegionNumericsFactory) Build() sharedregion.SharedRegionNumerics {
    desc := factory.info.UserDescription
    switch factory.info.DetectedNumericsMode {
    case NativeNumericsMode:
        nativeBaseApp := CreateNativeBaseFacade(factory.info, factory.baseApp)
        specialized := nativebase.CreateNativeBaseNumerics(nativeBaseApp)
        sequence := nativesequence.NewNativeSequenceNumerics(specialized)
        region := nativeregion.NewRegionNumerics(specialized, factory.regionConfig, sequence)
        return nativesharedregion.CreateNativeSharedRegion(region, factory.sharedConfig.Jobs)
    case BigFloatNumericsMode:
        bigBaseApp := CreateBigBaseFacade(factory.info, factory.baseApp)
        specialized := bigbase.Make(bigBaseApp)
        sequence := bigsequence.NewBigSequenceNumerics(specialized)
        region := bigregion.NewRegionNumerics(specialized, factory.regionConfig, sequence)
        return bigsharedregion.CreateBigSharedRegion(region, factory.sharedConfig.Jobs)
    default:
        log.Panic("Unknown numerics mode", factory.info.DetectedNumericsMode)
    }
}