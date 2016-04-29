package rest

import (
    "sync"
)

func ReadM(f func (), mut sync.RWMutex) {
    mut.RLock()
    defer mut.RUnlock()
    f()
}

func WriteM(f func (), mut sync.RWMutex) {
    mut.Lock()
    defer mut.Unlock()
    f()
}