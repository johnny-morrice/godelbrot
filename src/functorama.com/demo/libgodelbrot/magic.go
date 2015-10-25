package libgodelbrot

// See base/magic for a discussion of magic numbers as they pertain to this project

// MATHEMATICAL VALUES

// Named bignums
var bigZero big.Float = bigbase.CreateBigFloat(0, DefaultHighPrec)
var bigOne big.Float = bigbase.CreateBigFloat(1, DefaultHighPrec)
var bigTwo big.Float = bigbase.CreateBigFloat(2, DefaultHighPrec)

// DEFAULTS

// Default high precision for newly created big floats
const DefaultHighPrec uint = 500

const DefaultIterations uint8 = 255
const DefaultDivergeLimit float64 = 4.0
const DefaultImageWidth uint = 800
const DefaultImageHeight uint = 600
const DefaultCollapse uint = 8
const DefaultBufferSize uint = 256

// Default base for newly parsed numbers
const DefaultBase uint = 10

// What we consider to be a tiny image area, by default
const DefaultTinyImageArea uint = 40000

// What we consider to be a small number of render jobs
const DefaultLowThreading uint = 2

// Default sample size for region glitch-correction
const DefaultGlitchSamples uint = 10