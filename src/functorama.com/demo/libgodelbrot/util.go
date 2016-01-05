package libgodelbrot

import (
    "fmt"
    "errors"
    "log"
    "runtime"
    "math/big"
)

// I thought this might be a great idea,
// but it is horrendeous and must be removed.
func panic2err(factory func() interface{}) (anything interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            switch r := r.(type) {
            case runtime.Error:
                log.Panic(r)
            case error:
                err = r
            case string:
                err = errors.New(r)
            default:
                err = errors.New(fmt.Sprintf("%v", r))
            }
        }
    }()

    // Not obvious at first what this function would return if anything was uninitialized
    anything = nil
    anything = factory()

    return
}

// Panic to escape parsing of user input
func parsePanic(err error, inputName string) {
    log.Panic("Could not parse ", inputName, ":", err)
}

// Parse a big.Float
func parseBig(number string) (*big.Float, error) {
    f, _, err := big.ParseFloat(number, DefaultBase, DefaultHighPrec, big.ToNearestEven)
    return f, err
}