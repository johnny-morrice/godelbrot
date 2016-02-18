package draw

import (
	"image/color"
	"testing"
	"github.com/johnny-morrice/godelbrot/base"
)

func TestRedscaleColorBlack255(t *testing.T) {
	redscale := NewRedscalePalette(255)
	dark := base.BaseMandelbrot{InSet: true}
	expected := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	actual := redscale.Color(dark)
	if expected != actual {
		t.Error("Expected: ", expected, " Actual: ", actual)
	}
}

func TestRedscaleColorBright255(t *testing.T) {
	redscale := NewRedscalePalette(255)
	bright := base.BaseMandelbrot{InvDivergence: 0, InSet: false}
	expected := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	actual := redscale.Color(bright)
	if expected != actual {
		t.Error("Expected: ", expected, " Actual: ", actual)
	}
}

func TestRedscaleColorMid255(t *testing.T) {
	redscale := NewRedscalePalette(255)
	half := base.BaseMandelbrot{InvDivergence: 125, InSet: false}
	expected := color.NRGBA{R: 130, G: 0, B: 0, A: 255}
	actual := redscale.Color(half)
	if expected != actual {
		t.Error("Expected: ", expected, " Actual: ", actual)
	}
}

func TestRedscaleColorBright10(t *testing.T) {
	redscale := NewRedscalePalette(10)
	bright := base.BaseMandelbrot{InvDivergence: 0, InSet: false}
	expected := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	actual := redscale.Color(bright)
	if expected != actual {
		t.Error("Expected: ", expected, " Actual: ", actual)
	}
}

func TestRedscaleColorMid10(t *testing.T) {
	redscale := NewRedscalePalette(10)
	half := base.BaseMandelbrot{InvDivergence: 5, InSet: false}
	expected := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	actual := redscale.Color(half)
	if expected != actual {
		t.Error("Expected: ", expected, " Actual: ", actual)
	}
}
