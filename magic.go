package godelbrot

import (
	"github.com/johnny-morrice/godelbrot/internal/bigbase"
	"math/big"
)

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

// Minimum bounds of Mandelbrot set
const MandelbrotMin complex128 = -2.01 - 1.11i

// Maximum bounds of Mandelbrot set
const MandelbrotMax complex128 = 0.59 + 1.13i

// Named bignums
var bigZero big.Float = bigbase.MakeBigFloat(0, DefaultHighPrec)
var bigOne big.Float = bigbase.MakeBigFloat(1, DefaultHighPrec)
var bigTwo big.Float = bigbase.MakeBigFloat(2, DefaultHighPrec)

// DEFAULTS

// Default high precision for newly created big floats
const DefaultHighPrec uint = 500

// Default precision is for native arithmetic
const DefaultPrecision uint = 53

const DefaultIterations uint8 = 255
const DefaultDivergeLimit float64 = 4.0
const DefaultImageWidth uint = 600
const DefaultImageHeight uint = 600
const DefaultCollapse uint = 4
const DefaultBufferSize uint = 256

// Default base for newly parsed numbers
const DefaultBase int = 10

// What we consider to be a tiny image area, by default
const DefaultTinyImageArea uint = 40000

// Default sample size for region glitch-correction
const DefaultRegionSamples uint = 12
