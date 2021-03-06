package draw

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"image/color"
	"testing"
)

func TestColor(t *testing.T) {
	const iterLimit uint8 = 10
	cacher := func(iterLimit, index uint8) color.NRGBA {
		return color.NRGBA{index, index, index, 255}
	}
	white := color.NRGBA{255, 255, 255, 255}
	palette := NewCachePalette(iterLimit, white, cacher)

	inSet := base.EscapeValue{InSet: true}

	actualInSet := palette.Color(inSet)
	if white != actualInSet {
		t.Error("Expected white, but set member was assigned color:", actualInSet)
	}

	for i := uint8(0); i < iterLimit; i++ {
		expect := color.NRGBA{i, i, i, 255}
		member := base.EscapeValue{InvDiv: i}
		actual := palette.Color(member)

		if expect != actual {
			t.Error("Expected", expect, "but set member was assigned color:", actual)
		}
	}
}
