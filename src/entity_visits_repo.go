package main

import (
	// "fmt"
	"sync"
)

type EntityVisitsRepo struct {
	Collection    []*VisitsMRepo
	Mutex         sync.RWMutex
	MapCollection map[uint32]*VisitsMRepo
	MapMutex      sync.RWMutex
}

func makeEntityVisitsRepo(lenth uint32) EntityVisitsRepo {
	return EntityVisitsRepo{
		Collection:    make([]*VisitsMRepo, lenth),
		Mutex:         sync.RWMutex{},
		MapCollection: make(map[uint32]*VisitsMRepo),
		MapMutex:      sync.RWMutex{}}
}

func (repo *EntityVisitsRepo) findVisitsRepo(entity_id uint32) *VisitsMRepo {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	// var entity, _ = repo.Collection[entity_id]
	// return entity
	return repo.Collection[entity_id]
}

func (repo *EntityVisitsRepo) findOrInitVisitsRepo(entity_id uint32) *VisitsMRepo {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()
	visits_repo := repo.Collection[entity_id]
	if visits_repo == nil {
		visits_repo = &VisitsMRepo{
			Collection: make(map[uint32]*Visit),
			Mutex:      sync.RWMutex{}}
		repo.Collection[entity_id] = visits_repo
	}
	return visits_repo
}

func (repo *EntityVisitsRepo) addVisit(entity_id uint32, visit *Visit) {
	visits_repo := repo.findVisitsRepo(entity_id)
	if visits_repo == nil {
		visits_repo = repo.findOrInitVisitsRepo(entity_id)
	}
	visits_repo.Add(visit)
}
