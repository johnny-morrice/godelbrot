package libgodelbrot

// This module defines magic numbers

// Normalized size of box containing Mandelbrot set
const MagicSetSize complex128 = 2.1 + 2i

// Default offset for top left of plane containing set
const MagicOffset complex128 = -1.5 + 1i

// Defaults for the rendering system
// See RengerConfig
const DefaultIterations uint8 = 255
const DefaultDivergeLimit float64 = 4.0
const DefaultImageWidth uint = 800
const DefaultImageHeight uint = 600
const DefaultZoom float64 = 1.0
const DefaultCollapse uint  = 2