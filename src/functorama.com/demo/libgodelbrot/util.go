package libgodelbrot

func panic2err(factory func() interface{}) (anything interface{}, err) {
    defer func() {
        if r := recover(); r != nil {
            switch r := r.(type) {
            case runtime.Error:
                log.Panic(r)
            default:
                err = r.(error)
            }
        }
    }()

    // Not obvious at first what this function would return if anything was uninitialized
    anything := nil
    anything = factory()

    return
}

// Panic to escape parsing of user input
func parsePanic(err error, inputName string) {
    return log.Panic("Could not parse", inputName, ":", err)
}

// Parse a big.Float
func parseBig(number string) {
    // Do we need to care about the actual base used?
    f, _, err := big.ParseFloat(number, DefaultBase, DefaultHighPrec)
    return f, err
}