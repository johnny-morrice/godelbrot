package base

// This module covers internal details of the numeric strategies used by
// Godelbrot.
// The internal detail of each strategy is defined by an interfac.
// This is because the render control algorithms are also strategies that vary
// independently.

// PixelMember is a MandelbrotMember associated with a pixel
type PixelMember struct {
	I      int
	J      int
	Member MandelbrotMember
}
