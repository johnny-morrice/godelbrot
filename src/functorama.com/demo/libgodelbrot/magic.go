package libgodelbrot

// This module defines magic numbers relating to the Mandelbrot set

// Normalized size of window onto complex plane containing set
const MagicWindowSize complex128 = 2.1 + 2i

// Default offset for top left of plane contraining set
const MagicOffset complex128 = -1.5 + 1i