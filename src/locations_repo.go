package main

import (
	"strconv"
	"sync"
)

var (
	LOCATIONS_REPO_COLLECTION_SIZE = uint32(1010000)
)

type LocationsRepo struct {
	Collection []*Location
	// Mutex         sync.RWMutex
	MapCollection map[uint32]*Location
	MapMutex      sync.RWMutex
}

func (repo *LocationsRepo) InitEntity() *Location {
	entity := Location{}
	return &entity
}

func (repo *LocationsRepo) Create(data []byte) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
}

func (repo *LocationsRepo) Add(entity *Location) {
	// if entity.Id < LOCATIONS_REPO_COLLECTION_SIZE {
	repo.Collection[entity.Id] = entity
	// } else {
	// 	repo.MapMutex.Lock()
	// 	defer repo.MapMutex.Unlock()
	// 	repo.MapCollection[entity.Id] = entity
	// }
}

func (repo *LocationsRepo) Find(id uint32) (*Location, bool) {
	// if id < LOCATIONS_REPO_COLLECTION_SIZE {
	entity := repo.Collection[id]
	return entity, (entity != nil)
	// } else {
	// 	repo.MapMutex.Lock()
	// 	defer repo.MapMutex.Unlock()
	// 	entity, ok := repo.MapCollection[id]
	// 	return entity, ok
	// }
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
