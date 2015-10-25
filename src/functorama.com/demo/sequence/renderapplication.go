package sequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

type RenderApplication interface {
	base.RenderApplication
	draw.ContextProvider
	SequenceNumericsFactory() SequenceNumericsFactory
}