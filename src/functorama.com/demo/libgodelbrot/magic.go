package libgodelbrot

// This module defines magic numbers

// Normalized size of box containing Mandelbrot set
const MagicSetSize complex128 = 2.60 + 2.24i

// Default offset for top left of plane containing set
const MagicOffset complex128 = -2.01 + 1.13i

// Defaults for the rendering system
// See RengerConfig
const DefaultIterations uint8 = 255
const DefaultDivergeLimit float64 = 4.0
const DefaultImageWidth uint = 800
const DefaultImageHeight uint = 600
const DefaultZoom float64 = 1.0
const DefaultCollapse uint  = 2
const DefaultBufferSize uint = 256

// A fairly large number
const Meg uint = 1048576
// A smaller number
const Kilo uint = 1024
const K64 uint = Kilo * 64