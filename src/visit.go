package main

import (
	//"github.com/pquerna/ffjson/ffjson"
	"encoding/json"
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

func (entity *Visit) Update(data *JsonData, lock bool) {
	var location_id int
	var user_id int
	if lock {
		entity.Mutex.Lock()
	}
	for key, value := range *data {
		switch key {
		case "id":
			entity.Id = int(value.(float64))
		case "location":
			location_id = int(value.(float64))
			entity.LocationId = location_id
			location, _ := Locations.Find(location_id)
			entity.Location = location
			LocationsVisits.addVisit(location_id, entity)
		case "user":
			user_id = int(value.(float64))
			entity.UserId = user_id
			sync_user = true
			user, _ := Users.Find(user_id)
			entity.User = user
			UsersVisits.addVisit(user_id, entity)
		case "visited_at":
			entity.VisitedAt = int(value.(float64))
		case "mark":
			entity.Mark = int(value.(float64))
		}
	}
	if lock {
		entity.Mutex.Unlock()
	}
}

func (entity *Visit) toJson(w io.Writer) {
	entity.Mutex.RLock()
	json.NewEncoder(w).Encode(entity)
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

func (repo *VisitsRepo) Create(data *JsonData) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
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
