package rest

import (
    "fmt"
    "reflect"
    "sync"
)

type hashcode string

const emptyhash = hashcode("")

type hasher interface {
    hash() hashcode
}

type cache struct {
    cacheT reflect.Type
    items map[hashcode]hasher
    mutex sync.RWMutex
}

func makeCache(any interface{}) cache {
    ca := cache{}
    ca.cacheT = reflect.TypeOf(any)
    ca.items = map[hashcode]hasher{}
    return ca
}

func (ca *cache) put(h hasher) {
    if ht := reflect.TypeOf(h); ht != ca.cacheT {
        msg := fmt.Sprintf("Error: expected type %v but received type %v", ca.cacheT, ht)
        panic(msg)
    }

    hsh := h.hash()
    ca.mutex.Lock()
    defer ca.mutex.Unlock()
    ca.items[hsh] = h
}

func (ca *cache) get(hsh hashcode) (hasher, bool) {
    ca.mutex.RLock()
    defer ca.mutex.RUnlock()
    item, ok := ca.items[hsh]
    return item, ok
}