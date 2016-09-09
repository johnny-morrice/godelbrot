package draw

type ContextProvider interface {
	DrawingContext() DrawingContext
}
