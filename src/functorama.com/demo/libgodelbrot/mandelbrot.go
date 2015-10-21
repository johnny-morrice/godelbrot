package libgodelbrot

type MandelbrotMember interface {
	InverseDivergence() uint8
	SetMember() bool
	Mandelbrot(iterateLimit uint8, divergeLimit float64)
}

type BaseMandelbrot struct {
	InvDivergence uint8
	InSet         bool
}

func (base BaseMandelbrot) InverseDivergence() uint8 {
	return base.InvDivergence
}

func (base BaseMandelbrot) SetMember() bool {
	return base.InSet
}
