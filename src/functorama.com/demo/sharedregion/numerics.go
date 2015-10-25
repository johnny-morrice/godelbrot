package libgodelbrot

// SharedSequentialNumerics provides sequential (column-wise) rendering calculations for a threaded
// render strategy
type SharedSequentialNumerics interface {
    SequentialNumerics
    OpaqueThreadPrototype
}

// SharedRegionNumerics provides a RegionNumerics for threaded render stregies
type SharedRegionNumerics interface {
    RegionNumerics
    OpaqueThreadPrototype
}