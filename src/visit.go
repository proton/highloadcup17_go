package main

import (
	"encoding/json"
	// "fmt"
	"io"
	"sync"
)

type Visit struct {
	Id               uint32       `json:"id"`
	LocationId       uint32       `json:"location"`
	UserId           uint32       `json:"user"`
	VisitedAt        uint32       `json:"visited_at"`
	Mark             uint32       `json:"mark"`
	Mutex            sync.RWMutex `json:"-"`
	LocationPlace    string       `json:"-"`
	LocationCountry  string       `json:"-"`
	LocationDistance uint32       `json:"-"`
	UserBirthDate    int32        `json:"-"`
	UserGender       string       `json:"-"`
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

func (entity *Visit) Update(data *JsonData, lock bool) {
	sync_user := false
	sync_location := false
	if lock {
		entity.Mutex.Lock()
	}
	for key, value := range *data {
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
		location, _ := Locations.Find(entity.LocationId)
		if lock {
			location.Mutex.Lock()
		}
		entity.LocationPlace = location.Place
		entity.LocationCountry = location.Country
		entity.LocationDistance = location.Distance
		location.VisitIdsMap[entity.Id] = true
		if lock {
			location.Mutex.Unlock()
		}
	}

	if sync_user {
		// if entity.UserId == 84 {
		// 	fmt.Println("JsonData:", data)
		// }
		user, _ := Users.Find(entity.UserId)
		if lock {
			user.Mutex.Lock()
		}
		entity.UserBirthDate = user.BirthDate
		entity.UserGender = user.Gender
		user.VisitIdsMap[entity.Id] = true
		if lock {
			user.Mutex.Unlock()
		}
	}

	if lock {
		entity.Mutex.Unlock()
	}
}

func (entity *Visit) to_json(w io.Writer) {
	entity.Mutex.RLock()
	json.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

func (visit *Visit) ToView() *VisitView {
	return &VisitView{
		Mark:      visit.Mark,
		VisitedAt: visit.VisitedAt,
		Place:     visit.LocationPlace}
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

func (repo *VisitsRepo) Create(data *JsonData) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
}

func (repo *VisitsRepo) CreateFromJson(raw_data []byte) error {
	entity := repo.InitEntity()
	err := json.Unmarshal(raw_data, entity)
	if err == nil {
		repo.Add(entity)
	}
	return err
}

func (repo *VisitsRepo) Add(entity *Visit) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *VisitsRepo) Find(id uint32) (*Visit, bool) {
	repo.Mutex.RLock()
	var entity, found = repo.Collection[id]
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

func (repo *VisitsRepo) FindEntity(id uint32) (Entity, bool) {
	return repo.Find(id)
}
