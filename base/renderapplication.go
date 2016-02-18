package base

// A facade used by subsystems to interact with the application at large
type RenderApplication interface {
	// Basic configuration
	BaseConfig() BaseConfig
	PictureDimensions() (uint, uint)
}

type BaseConfig struct {
	IterateLimit uint8
	DivergeLimit float64
	FixAspect bool
}
