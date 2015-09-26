package libgodelbrot


type RedscalePalette CachePalette

func NewRedscalePalette(iterateLimit uint8) RedscalePalette {
    return NewCachePalette(iterateLimit, &redscaleCacher)
}

// Cache redscale colour values
func redscaleCacher(limit uint8, index uint8) color.NRGBA {
    return  color.NRGBA{
        R: limit - index,
        G: 0,
        B: 0,
        A: 255,
    }
}