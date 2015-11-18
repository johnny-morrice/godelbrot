package sequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

type SequenceNumericsFactory interface {
    Build() SequenceNumerics
}

type SequenceProvider interface {
    SequenceNumericsFactory() SequenceNumericsFactory
}

type RenderApplication interface {
	base.RenderApplication
	draw.ContextProvider
	SequenceProvider
}