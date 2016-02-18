package draw

type MockContextProvider struct {
    TDrawingContext bool

    Context DrawingContext
}

func (mock *MockContextProvider) DrawingContext() DrawingContext {
    mock.TDrawingContext = true
    return mock.Context
}