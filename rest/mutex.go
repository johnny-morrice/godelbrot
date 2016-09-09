package rest

import (
	"sync"
)

func readM(mut sync.RWMutex, f func()) {
	mut.RLock()
	defer mut.RUnlock()
	f()
}

func writeM(mut sync.RWMutex, f func()) {
	mut.Lock()
	defer mut.Unlock()
	f()
}

func readwriteM(mut sync.Mutex, f func()) {
	mut.Lock()
	defer mut.Unlock()
	f()
}
