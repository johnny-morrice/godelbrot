package base

func PictureAspectRatio(width uint, height uint) float64 {
	return float64(width) / float64(height)
}

func AppPictureAspectRatio(app RenderApplication) float64 {
	return PictureAspectRatio(app.PictureDimensions())
}
