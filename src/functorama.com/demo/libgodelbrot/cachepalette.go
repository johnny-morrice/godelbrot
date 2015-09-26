package libgodelbrot

type Cacher func (iterateLimit uint8, index uint8) color.NRGBA

type CachePalette struct {
    setColor color.NRGBA
    scale []color.NRGBA
    limit uint8
}

func NewCachePalette(iterateLimit uint8, cacher Cacher) CachePalette {
    palette := make([]color.NRGBA, iterateLimit, iterateLimit)
    for i := 0; i < iterateLimit; i++ {
        palette.scale[i] = cacher.cache(iteraterLimit, uint8(i))
    }
    return palette
}

// CachePalette implements Palette
func (palette CachePalette) Color(member MandelbrotMember) color.NRGBA {
    if member.InSet {
        return palette.setColor
        }
    } else {
        return palette.scale[member.InvDivergence]
    }
}