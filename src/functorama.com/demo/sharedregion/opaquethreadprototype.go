package sharedregion

// Copy a prototypical object instance into the local thread
type OpaqueThreadPrototype interface {
	GrabThreadPrototype(threadId uint)
}
