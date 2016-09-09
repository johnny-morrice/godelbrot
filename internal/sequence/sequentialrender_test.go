package sequence

import (
	"github.com/johnny-morrice/godelbrot/internal/draw"
	"image"
	"testing"
)

func TestMake(t *testing.T) {
	expectedPic := image.NewNRGBA(image.ZR)
	context := &draw.MockDrawingContext{
		Pic: expectedPic,
	}
	expectedNumerics := &MockNumerics{}
	expectedRenderer := SequenceRenderStrategy{
		numerics: expectedNumerics,
		context:  context,
	}
	factory := &MockFactory{Numerics: expectedNumerics}
	mock := &MockRenderApplication{
		SequenceFactory: factory,
	}
	mock.Context = context
	mock.Base.IterateLimit = 200
	actualRenderer := Make(mock)

	if !(mock.TSequenceNumericsFactory && mock.TDrawingContext) {
		t.Error("Expected methods not called on mock:", mock)
	}

	if !factory.TBuild {
		t.Error("Expected methods not called on mock sequence factory:", factory)
	}

	if actualRenderer != expectedRenderer {
		t.Error("Expected renderer ", expectedRenderer,
			"not equal to actual renderer:", actualRenderer)
	}
}

func TestStrategyRender(t *testing.T) {
	const ilimit = 255
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	expectedPic := image.NewNRGBA(image.ZR)
	context := draw.NewMockDrawingContext(ilimit)
	context.Pic = expectedPic
	numerics := &MockNumerics{}

	renderer := SequenceRenderStrategy{
		numerics: numerics,
		context:  context,
	}
	actualPic, err := renderer.Render()

	if err != nil {
		t.Error("Unexpected error in render")
	}

	if actualPic != expectedPic {
		t.Error("Expected a certain picture to be returned but was:", actualPic)
	}

	if !numerics.TSequence {
		t.Error("Expected methods not called on numerics:", numerics)
	}

	if !context.TPicture {
		t.Error("Expected methods not called on context:", context)
	}
}
