package godelbrot

import (
	"testing"

	"github.com/johnny-morrice/godelbrot/config"
)

func TestConfigure(t *testing.T) {
	info, noErr := Configure(DefaultRequest())

	if info == nil {
		t.Error("RenderContext created from default description should not be nil")
	}

	if noErr != nil {
		t.Error("There should be no error from constructing the default RenderContext")
	}

	noInfo, realErr := Configure(&config.Request{})

	if realErr == nil {
		t.Error("Expect error after trying to create context from no description")
	}

	if noInfo != nil {
		t.Error("Did not expect RenderContext after blank construction.")
	}
}