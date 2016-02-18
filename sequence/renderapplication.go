package sequence

import (
	"github.com/johnny-morrice/godelbrot/base"
	"github.com/johnny-morrice/godelbrot/draw"
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