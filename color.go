package godelbrot

import (
	"github.com/johnny-morrice/godelbrot/internal/base"
	"image"
)

func Recolor(desc *Info, gray image.Image) *image.NRGBA {
	// CAUTION lossy conversion
	iterlim := desc.UserRequest.IterateLimit

	dfac := makeDrawFacade(desc)
	palette := dfac.Colors()

	bnd := gray.Bounds()
	bright := image.NewNRGBA(bnd)

	for x := bnd.Min.X; x < bnd.Max.X; x++ {
		for y := bnd.Min.Y; y < bnd.Max.Y; y++ {
			bigdiv, _, _, _ := gray.At(x, y).RGBA()
			invdiv := uint8(bigdiv >> 8)
			member := base.EscapeValue{
				InvDiv: invdiv,
				InSet:  invdiv == iterlim,
			}
			col := palette.Color(member)
			bright.Set(x, y, col)
		}
	}

	return bright
}
