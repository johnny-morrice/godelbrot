package libgodelbrot

// Draw the Mandelbrot set.  This is the main entry point to libgodelbrot
func Godelbrot(desc RenderDescription) (image.NRGBA, error) {
    facade, err := GodelbrotRenderContext(desc)
    if err != nil {
        return nil, error
    }

    return facade.Render(), nil
}

// Based on the description, choose a renderer, numerical system and palette
// and combine them into a coherent render context
func GodelbrotRenderContext(desc RenderDescription) (RenderContext, error) {
    initializer, err := InitializeContext(desc)
    if err != nil {
        return nil, err
    }
    return initializer.CreateFacade(), nil
}