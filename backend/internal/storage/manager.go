package storage

import "sync"

var (
    current Storage
    mu sync.RWMutex
)

func Set(s Storage) { mu.Lock(); current = s; mu.Unlock() }
func Get() Storage { mu.RLock(); defer mu.RUnlock(); return current }