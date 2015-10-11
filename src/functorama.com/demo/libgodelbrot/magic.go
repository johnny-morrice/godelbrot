package libgodelbrot

import (
    "math/big"
)

// This module defines magic numbers.

// There are three different kinds of magic number in this program:
// 1. Mathematical constants.  For example, the location of the Mandelbrot set 
//    is unlikely to change any time soon.  These are the most truly magical 
//    numbers in this program.
// 2. Defaults.  It is convenient for this library to provide some sensible
//    defaults for its parameters.  As these are "guaranteed to work", these 
//    represent mathematical truths in their own way.
// 3. Operational constants.  This program uses various integers in control
//    flow, in arithmetic, and when requesting memory from the operating system.
//    These are the least "magical" numbers present in this program.  Each of 
//    these should be reviewed, since it may be desirable to replace these with 
//    an item that can be configured at runtime.

// MATHEMATICAL CONSTANTS

// Normalized size of box containing Mandelbrot set
const MagicSetSize complex128 = 2.60 + 2.24i
// Default offset for top left of plane containing set
const MagicOffset complex128 = -2.01 + 1.13i
// Named bignums
var bigZero big.Float = NewFloat(0)
var bigOne big.Float = NewFloat(1)
var bigTwo big.Float = NewFloat(2)

// DEFAULTS

// Defaults for the rendering system
// See RengerConfig
const DefaultIterations uint8 = 255
const DefaultDivergeLimit float64 = 4.0
const DefaultImageWidth uint = 800
const DefaultImageHeight uint = 600
const DefaultCollapse uint = 2
const DefaultBufferSize uint = 256
// Default base for newly parsed numbers
const DefaultBase uint = 10
// Default high precision for newly created big floats
const DefaultHighPrec uint = 500
// What we consider to be a tiny image area, by default
const DefaultTinyImageArea uint = 40000
// What we consider to be a small number of render jobs
const DefaultLowThreading uint = 2

// OPERATIONAL CONSTANTS

// A medium sized number
const allocMedium uint = allocSmall * 64
// A small number
const allocSmall uint = 1024