package libgodelbrot


type Palette interface {
    Color(point MandelbrotMember)
}