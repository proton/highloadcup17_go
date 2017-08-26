package main

import (
	"github.com/buger/jsonparser"
	"github.com/pquerna/ffjson/ffjson"
	// "encoding/json"
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
	Json       []byte       `json:"-"`
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
	sync_user := false
	sync_location := false
	if lock {
		entity.Mutex.Lock()
	}
	for key, value := range *data {
		switch key {
		case "id":
			entity.Id = int(value.(float64))
		case "location":
			entity.LocationId = int(value.(float64))
			sync_location = true
			location, _ := Locations.Find(entity.LocationId)
			entity.Location = location
		case "user":
			entity.UserId = int(value.(float64))
			sync_user = true
			user, _ := Users.Find(entity.UserId)
			entity.User = user
		case "visited_at":
			entity.VisitedAt = int(value.(float64))
		case "mark":
			entity.Mark = int(value.(float64))
		}
	}
	entity.cacheJSON()
	if lock {
		entity.Mutex.Unlock()
	}

	if sync_location {
		LocationsVisits.addVisit(entity.LocationId, entity)
	}
	if sync_user {
		UsersVisits.addVisit(entity.UserId, entity)
	}
}

func (entity *Visit) UpdateFromJSON(data []byte, lock bool) {
	if lock {
		entity.Mutex.Lock()
	}
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// fmt.Println(jsonparser.Get(value, "url"))
	}, "person", "avatars")
	// for key, value := range *data {
	// 	switch key {
	// 	case "id":
	// 		entity.Id = int(value.(float64))
	// 	case "place":
	// 		entity.Place = value.(string)
	// 	case "country":
	// 		entity.Country = value.(string)
	// 	case "city":
	// 		entity.City = value.(string)
	// 	case "distance":
	// 		entity.Distance = int(value.(float64))
	// 	}
	// }
	entity.cacheJSON()
	if lock {
		entity.Mutex.Unlock()
	}
}

func (entity *Visit) cacheJSON() {
	b, _ := ffjson.Marshal(entity)
	entity.Json = b
}

func (entity *Visit) writeJSON(w io.Writer) {
	entity.Mutex.RLock()
	w.Write(entity.Json)
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

func (repo *VisitsRepo) CreateFromJSON(data []byte) {
	entity := repo.InitEntity()
	entity.UpdateFromJSON(data, false)
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
