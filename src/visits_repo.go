package main

import (
	"sync"
)

var (
	VISITS_REPO_COLLECTION_SIZE = uint32(10100000)
)

type VisitsRepo struct {
	Collection    []*Visit
	MapCollection map[uint32]*Visit
	MapMutex      sync.RWMutex
}

func makeVisitsRepo(lenth uint32) VisitsRepo {
	return VisitsRepo{
		Collection:    make([]*Visit, lenth),
		MapCollection: make(map[uint32]*Visit),
		MapMutex:      sync.RWMutex{}}
}

func (repo *VisitsRepo) InitEntity() *Visit {
	entity := Visit{}
	return &entity
}

func (repo *VisitsRepo) Create(data []byte) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
}

func (repo *VisitsRepo) Add(entity *Visit) {
	if entity.Id < VISITS_REPO_COLLECTION_SIZE {
		repo.Collection[entity.Id] = entity
	} else {
		// repo.MapMutex.Lock()
		// defer repo.MapMutex.Unlock()
		repo.MapCollection[entity.Id] = entity
	}
}

func (repo *VisitsRepo) Find(id uint32) (*Visit, bool) {
	if id < VISITS_REPO_COLLECTION_SIZE {
		entity := repo.Collection[id]
		return entity, (entity != nil)
	} else {
		// 	repo.MapMutex.Lock()
		// 	defer repo.MapMutex.Unlock()
		entity, ok := repo.MapCollection[id]
		return entity, ok
	}
}

func (repo *VisitsRepo) FindEntity(id uint32) (Entity, bool) {
	return repo.Find(id)
}
