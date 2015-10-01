package libgodelbrot

// This module defines magic numbers relating to the Mandelbrot set

// Normalized size of box containing Mandelbrot set
const MagicSetSize complex128 = 2.1 + 2i

// Default offset for top left of plane containing set
const MagicOffset complex128 = -1.5 + 1i