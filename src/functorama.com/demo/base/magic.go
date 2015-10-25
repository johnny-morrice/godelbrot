package base

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

// Minimum bounds of Mandelbrot set
const MandelbrotMin complex128 = -2.01 - 1.11i

// Maximum bounds of Mandelbrot set
const MandelbrotMax complex128 = 0.59 + 1.13i

// OPERATIONAL CONSTANTS

// A large number
const AllocLarge uint = AllocMedium * 1024

// A medium sized number
const AllocMedium uint = AllocSmall * 64

// A small number
const AllocSmall uint = 1024

// A tiny allocation
const AllocTiny uint = 128