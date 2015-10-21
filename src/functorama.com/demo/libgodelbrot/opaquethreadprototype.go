package libgodelbrot

// Copy a prototypical object instance into the local thread
type OpaqueThreadedPrototype interface {
	GrabThreadPrototype(threadId uint)
}
