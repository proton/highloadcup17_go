package main

import (
	// "fmt"
	"sync"
)

type EntityVisitsRepo struct {
	Collection map[int]*VisitsRepo
	Mutex      sync.RWMutex
}

func (repo *EntityVisitsRepo) findVisitsRepo(entity_id int) *VisitsRepo {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	var entity, _ = repo.Collection[entity_id]
	return entity
}

func (repo *EntityVisitsRepo) initVisitsRepo(entity_id int) *VisitsRepo {
	repo.Mutex.Lock()
	visits_repo := repo.Collection[entity_id]
	if visits_repo == nil {
		visits_repo = &VisitsRepo{
			Collection: make(map[int]*Visit),
			Mutex:      sync.RWMutex{}}
	}
	repo.Collection[entity_id] = visits_repo
	repo.Mutex.Unlock()
	return visits_repo
}

func (repo *EntityVisitsRepo) addVisit(entity_id int, visit *Visit) {
	visits_repo := repo.findVisitsRepo(entity_id)
	if visits_repo == nil {
		visits_repo = repo.initVisitsRepo(entity_id)
	}
	visits_repo.Add(visit)
}
