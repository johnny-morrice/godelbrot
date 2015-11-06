package libgodelbrot

type SequenceFacade struct {
    *BaseFacade
    *DrawFacade
    factory SequenceNumericsFactory
}

var _ sequence.RenderApplication = (*SequenceFacade)(nil)

func NewSequenceFacade(info *RenderInfo) *SequenceFacade {
    baseApp := NewBaseFacade(info)
    facade := &SequenceFacade{
        BaseFacade: baseApp,
        DrawFacade: NewDrawFacade(info),
    }
    facade.factory = &GodelbrotSequenceNumericsFactory{info, baseApp}
    return facade
}

func (facade *SequenceFacade) SequenceNumericsFactory() sequence.SequenceNumericsFactory {
    return facade.factory
}

type GodelbrotSequenceNumericsFactory struct {
    info RenderInfo
    baseApp base.RenderApplication
}

func (factory *GodelbrotSequenceNumericsFactory) Build() sequence.SequenceNumerics {
    switch factory.info.DetectedNumericsMode {
    case NativeNumericsMode:
        nativeBaseApp := CreateNativeBaseFacade(factory.info, factory.baseApp)
        specialized := nativebase.CreateNativeBaseNumerics(nativeBaseApp)
        return nativesequence.NewNativeSequenceNumerics(specialized)
    case BigFloatNumericsMode:
        bigBaseApp := CreateBigBaseFacade(factory.info, factory.baseApp)
        specialized := bigbase.CreateBigBaseNumerics(bigBaseApp)
        return bigsequence.NewBigSequenceNumerics(specialized)
    default:
        log.Panic("Unknown numerics mode", factory.info.DetectedNumericsMode)
    }
}