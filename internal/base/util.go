package base

import (
	"fmt"
	"os"
)

func Dbg(things ...interface{}) {
	fmt.Fprintln(os.Stderr, things...)
}
