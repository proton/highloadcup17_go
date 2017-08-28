package main

import (
	"sync"
)

type VisitsMRepo struct {
	Collection map[uint32]*Visit
	Mutex      sync.RWMutex
}

func (repo *VisitsMRepo) Add(entity *Visit) {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()
	repo.Collection[entity.Id] = entity
}
