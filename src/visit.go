package main

import (
	"github.com/pquerna/ffjson/ffjson"
	// "fmt"
	"io"
	"sync"
)

type Visit struct {
	Id         uint32       `json:"id"`
	LocationId uint32       `json:"location"`
	UserId     uint32       `json:"user"`
	VisitedAt  uint32       `json:"visited_at"`
	Mark       uint32       `json:"mark"`
	Mutex      sync.RWMutex `json:"-"`
	Location   *Location    `json:"-"`
	User       *User        `json:"-"`
	// LocationPlace    string       `json:"-"`
	// LocationCountry  string       `json:"-"`
	// LocationDistance uint32       `json:"-"`
	// UserBirthDate    int32        `json:"-"`
	// UserGender       string       `json:"-"`
}

type VisitView struct {
	Mark      uint32 `json:"mark"`
	VisitedAt uint32 `json:"visited_at"`
	Place     string `json:"place"`
}

type VisitsRepo struct {
	Collection map[uint32]*Visit
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
			entity.Id = uint32(value.(float64))
		case "location":
			entity.LocationId = uint32(value.(float64))
			sync_location = true
		case "user":
			entity.UserId = uint32(value.(float64))
			sync_user = true
		case "visited_at":
			entity.VisitedAt = uint32(value.(float64))
		case "mark":
			entity.Mark = uint32(value.(float64))
		}
	}

	if sync_location {
		location, _ := Locations.Find(entity.LocationId, lock)
		entity.Location = location
		// if lock {
		// 	location.Mutex.Lock()
		// }
		// entity.LocationPlace = location.Place
		// entity.LocationCountry = location.Country
		// entity.LocationDistance = location.Distance
		// location.VisitIdsMap[entity.Id] = true
		// if lock {
		// 	location.Mutex.Unlock()
		// }
	}

	if sync_user {
		user, _ := Users.Find(entity.UserId, lock)
		entity.User = user
		// if lock {
		// 	user.Mutex.Lock()
		// }
		// entity.UserBirthDate = user.BirthDate
		// entity.UserGender = user.Gender
		// user.VisitIdsMap[entity.Id] = true
		// if lock {
		// 	user.Mutex.Unlock()
		// }
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

func NewVisitsRepo() VisitsRepo {
	return VisitsRepo{
		Collection: make(map[uint32]*Visit),
		Mutex:      sync.RWMutex{}}
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

func (repo *VisitsRepo) Find(id uint32, lock bool) (*Visit, bool) {
	if lock {
		// repo.Mutex.RLock()
	}
	repo.Mutex.RLock()
	var entity, found = repo.Collection[id]
	if lock {
		// repo.Mutex.RUnlock()
	}
	repo.Mutex.RUnlock()
	return entity, found
}

func (repo *VisitsRepo) FindAll(ids []uint32) ([]*Visit, bool) {
	entities := make([]*Visit, len(ids))
	repo.Mutex.RLock()
	for i, id := range ids {
		var entity, _ = repo.Collection[id]
		// entity lock?
		entities[i] = entity
	}
	repo.Mutex.RUnlock()
	return entities, true
}

func (repo *VisitsRepo) FindEntity(id uint32, lock bool) (Entity, bool) {
	return repo.Find(id, lock)
}
