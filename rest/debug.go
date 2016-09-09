package rest

import (
	"log"
)

var __DEBUG bool = false

func Debug() {
	__DEBUG = true
}

func debugf(format string, many ...interface{}) {
	if __DEBUG {
		log.Printf(format, many...)
	}
}
