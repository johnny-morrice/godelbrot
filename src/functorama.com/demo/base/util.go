package base

import (
    "os"
    "fmt"
)

func Dbg(things... interface{}) {
    fmt.Fprintln(os.Stderr, things...)
}