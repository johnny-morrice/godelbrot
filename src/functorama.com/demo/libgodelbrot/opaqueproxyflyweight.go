package libgodelbrot

// OpaqueProxyFlyweight enables implementation of its eponymous design pattern. You need lots of
// objects, but creating them is expensive. You want to externalize state, but your object is opaque
// (you access it through an interface) Only these objects can know the type of the external state.
// So create an extension of this object that acts like a proxy When needed, it will grab the
// extrinsic, without exposing the details to the client.
type OpaqueProxyFlyweight interface {
    ClaimExtrinics()
}

