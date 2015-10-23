package libgodelbrot

import (
	"testing"
)

func TestGodelbrotRenderContext(t *testing.T) {
	app, noErr := GodelbrotRenderContext(DefaultRenderDescription())

	if app == nil {
		t.Error("RenderContext created from default description should not be nil")
	}

	if noErr != nil {
		t.Error("There should be no error from constructing the default RenderContext")
	}

	noApp, realErr := GodelbrotRenderContext(RenderDescription{})

	if realErr == nil {
		t.Error("Expect error after trying to create context from no description")
	}

	if noApp != nil {
		t.Error("Did not expect RenderContext after blank construction.")
	}
}
