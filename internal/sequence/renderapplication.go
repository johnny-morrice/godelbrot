package sequence

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"github.com/johnny-morrice/godelbrot/internal/draw"
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