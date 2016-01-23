package libgodelbrot

import (
    "image"
    "math"
    "functorama.com/demo/base"
)

func Recolor(desc *Info, gray image.Image) *image.NRGBA {
    // CAUTION lossy conversion
    iterlim := desc.UserRequest.IterateLimit

    dfac := makeDrawFacade(desc)
    palette := dfac.Colors()

    bnd := gray.Bounds()
    bright := image.NewNRGBA(bnd)
    scale := float64(0xff) / float64(0xffff)

    for x := bnd.Min.X; x < bnd.Max.X; x++ {
        for y := bnd.Min.Y; y < bnd.Max.Y; y++ {
            bigdiv, _, _, _ := gray.At(x, y).RGBA()
            invdiv := uint8(math.Floor(scale * float64(bigdiv)))
            member := base.BaseMandelbrot{
                InvDivergence: invdiv,
                InSet: invdiv == iterlim,
            }
            col := palette.Color(member)
            bright.Set(x, y, col)
        }
    }

    return bright
}