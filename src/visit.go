package main

import (
	"github.com/pquerna/ffjson/ffjson"
	// "fmt"
	"io"
	"sync"
)

type Visit struct {
	Id         int          `json:"id"`
	LocationId int          `json:"location"`
	UserId     int          `json:"user"`
	VisitedAt  int          `json:"visited_at"`
	Mark       int          `json:"mark"`
	Mutex      sync.RWMutex `json:"-"`
	Location   *Location    `json:"-"`
	User       *User        `json:"-"`
}

type VisitView struct {
	Mark      int    `json:"mark"`
	VisitedAt int    `json:"visited_at"`
	Place     string `json:"place"`
}

type VisitsRepo struct {
	Collection map[int]*Visit
	Mutex      sync.RWMutex
}

func (entity *Visit) Update(data *JsonData, lock bool) bool {
	sync_user := false
	sync_location := false
	if lock {
		entity.Mutex.Lock()
		defer entity.Mutex.Unlock()
	}
	for key, value := range *data {
		if value == nil {
			return false
		}
		switch key {
		case "id":
			entity.Id = int(value.(float64))
		case "location":
			entity.LocationId = int(value.(float64))
			sync_location = true
		case "user":
			entity.UserId = int(value.(float64))
			sync_user = true
		case "visited_at":
			entity.VisitedAt = int(value.(float64))
		case "mark":
			entity.Mark = int(value.(float64))
		}
	}

	if sync_location {
		LocationsVisits.addVisit(entity.LocationId, entity)
	}
	if sync_user {
		UsersVisits.addVisit(entity.UserId, entity)
	}

	return true
}

func (entity *Visit) to_json(w io.Writer) {
	entity.Mutex.RLock()
	ffjson.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

func (visit *Visit) ToView() *VisitView {
	return &VisitView{
		Mark:      visit.Mark,
		VisitedAt: visit.VisitedAt,
		Place:     visit.Location.Place}
}

func (repo *VisitsRepo) InitEntity() *Visit {
	entity := Visit{}
	return &entity
}

func (repo *VisitsRepo) Create(data *JsonData) bool {
	entity := repo.InitEntity()
	ok := entity.Update(data, false)
	if !ok {
		return false
	}
	repo.Add(entity)
	return true
}

func (repo *VisitsRepo) Add(entity *Visit) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *VisitsRepo) Find(id int) (*Visit, bool) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	entity, ok := repo.Collection[id]
	return entity, ok
}

func (repo *VisitsRepo) FindEntity(id int) (Entity, bool) {
	return repo.Find(id)
}

// func (repo *VisitsRepo) FindAll(ids []int) []*Visit {
// 	entities := make([]*Visit, len(ids))
// 	repo.Mutex.RLock()
// 	for i, id := range ids {
// 		var entity, _ = repo.Collection[id]
// 		// entity lock?
// 		entities[i] = entity
// 	}
// 	repo.Mutex.RUnlock()
// 	return entities
// }
