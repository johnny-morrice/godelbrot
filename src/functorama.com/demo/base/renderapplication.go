package base

// A facade used by subsystems to interact with the application at large
type RenderApplication interface {
	// Basic configuration
	IterateLimit() uint8
	DivergeLimit() float64
	PictureDimensions() (uint, uint)
	FixAspect() bool
}
