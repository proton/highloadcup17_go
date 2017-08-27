package main

import (
	"strconv"
	"sync"
)

type LocationsRepo struct {
	Collection map[uint32]*Location
	Mutex      sync.RWMutex
}

func (repo *LocationsRepo) InitEntity() *Location {
	return &Location{}
}

func (repo *LocationsRepo) Create(data []byte) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
}

func (repo *LocationsRepo) Add(entity *Location) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *LocationsRepo) Find(id uint32) (*Location, bool) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	entity, ok := repo.Collection[id]
	return entity, ok
}

func (repo *LocationsRepo) FindEntity(id uint32) (Entity, bool) {
	return repo.Find(id)
}

func find_location(entity_id_str []byte) (*Location, bool) {
	entity_id_int, error := strconv.Atoi(bstring(entity_id_str))
	if error == nil {
		entity_id := uint32(entity_id_int)
		return Locations.Find(entity_id)
	}
	return nil, false
}
