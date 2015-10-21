package libgodelbrot

import (
    "testing"
    "image/color"
)

func TestColor(t *testing.T) {
    iterLimit := 10
    cacher := func (iterLimit, index uint8) color.RGBA {
            gray := color.Gray(index)
            return color.RGBA{gray.RGBA}
        }
    }
    white := color.RGBA{255,255,255,255}
    palette := NewCachePalette(iterLimit, white, cacher)

    inSet := BaseMandelbrot{InSet: true}

    actualInSet := palette.Color(inSet)
    if white != actualInSet {
        t.Error("Expected white, but set member was assigned color:", actualInSet)
    }

    for i := 0; i < iterLimit; i++ {
        expect := color.RGBA{i, i, i, i}
        member := BaseMandelbrot{InvDivergence: i}
        actual := palette.Color(member)

        if expect != actual {
            t.Error("Expected", expect, "but set member was assigned color:", actual)
        }
    }
}