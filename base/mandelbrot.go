package base

type MandelbrotMember struct {
	InvDiv uint8
	InSet         bool
}

// PixelMember is a MandelbrotMember associated with a pixel
type PixelMember struct {
    I      int
    J      int
    Member MandelbrotMember
}