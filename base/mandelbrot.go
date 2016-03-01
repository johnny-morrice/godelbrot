package base

type EscapeValue struct {
	InvDiv uint8
	InSet         bool
}

// PixelMember is a EscapeValue associated with a pixel
type PixelMember struct {
    I      int
    J      int
    Member EscapeValue
}