package libgodelbrot

func PictureAspectRatio(width uint, height uint) {
    return width / height
}

func AppPictureAspectRatio(app RenderApplication) {
    return PictureAspectRatio(app.PictureDimensions())
}