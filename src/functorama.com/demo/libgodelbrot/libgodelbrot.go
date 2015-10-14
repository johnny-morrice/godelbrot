package libgodelbrot

// Draw the Mandelbrot set.  This is the main entry point to libgodelbrot
func Godelbrot(desc RenderDescription) (image.NRGBA, error) {
    facade, _, err := GodelbrotRenderContext(desc)
    if err != nil {
        return nil, error
    }

    return facade.Render(), nil
}

// Based on the description, choose a renderer, numerical system and palette
// and combine them into a coherent render context
// Return the context and information on its settings
func GodelbrotRenderContext(desc RenderDescription) (RenderContext, RenderInfo, error) {
    initializer, err := InitializeContext(desc)
    if err != nil {
        return nil, err
    }
    return initializer.NewUserFacade(), initializer.info, nil
}