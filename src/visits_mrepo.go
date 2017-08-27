package main

import (
	"sync"
)

type VisitsMRepo struct {
	Collection map[uint32]*Visit
	Mutex      sync.RWMutex
}

func (repo *VisitsMRepo) InitEntity() *Visit {
	entity := Visit{}
	return &entity
}

func (repo *VisitsMRepo) Create(data []byte) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
}

func (repo *VisitsMRepo) Add(entity *Visit) {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()
	repo.Collection[entity.Id] = entity
}

func (repo *VisitsMRepo) Find(id uint32) (*Visit, bool) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	entity, ok := repo.Collection[id]
	return entity, ok
}

func (repo *VisitsMRepo) FindEntity(id uint32) (Entity, bool) {
	return repo.Find(id)
}
