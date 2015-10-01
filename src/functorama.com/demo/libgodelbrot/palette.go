package libgodelbrot

import (
    "image/color"
)

type Palette interface {
    Color(point MandelbrotMember) color.NRGBA
}