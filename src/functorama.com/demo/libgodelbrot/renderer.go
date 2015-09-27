package libgodelbrot

type Renderer func (argP *RenderParameters, palette Palette) (*image.NRGBA, error)