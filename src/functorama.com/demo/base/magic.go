package base

// This module defines magic numbers.

// See 'libgodelbrot/magic.go' for a philosophy on types of constants

// OPERATIONAL CONSTANTS

// A large number
const AllocLarge uint = AllocMedium * 1024

// A medium sized number
const AllocMedium uint = AllocSmall * 64

// A small number
const AllocSmall uint = 1024

// A tiny allocation
const AllocTiny uint = 128