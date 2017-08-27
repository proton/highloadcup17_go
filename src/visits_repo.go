package main

import (
	"sync"
)

type VisitsRepo struct {
	Collection map[uint32]*Visit
	Mutex      sync.RWMutex
}

func (repo *VisitsRepo) InitEntity() *Visit {
	entity := Visit{}
	return &entity
}

func (repo *VisitsRepo) CreateFromJSON(data []byte) {
	entity := repo.InitEntity()
	entity.UpdateFromJSON(data, false)
	repo.Add(entity)
}

func (repo *VisitsRepo) Add(entity *Visit) {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()
	repo.Collection[entity.Id] = entity
}

func (repo *VisitsRepo) Find(id uint32) (*Visit, bool) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	entity, ok := repo.Collection[id]
	return entity, ok
}

func (repo *VisitsRepo) FindEntity(id uint32) (Entity, bool) {
	return repo.Find(id)
}
