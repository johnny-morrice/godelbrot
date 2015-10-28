package sequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

type SequenceNumericsFactory interface {
    Build() SequenceNumerics
}

type RenderApplication interface {
	base.RenderApplication
	draw.ContextProvider
	SequenceNumericsFactory() SequenceNumericsFactory
}