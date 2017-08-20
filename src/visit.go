package main

import (
	"encoding/json"
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
}

type VisitsRepo struct {
	Collection map[uint32]*Visit
	Mutex      sync.RWMutex
}

func (entity *Visit) Update(data *JsonData) {
	entity.Mutex.Lock()
	for key, value := range *data {
		switch key {
		case "id":
			entity.Id = uint32(value.(float64))
		case "location":
			entity.LocationId = uint32(value.(float64))
			// TODO:
		case "user":
			entity.UserId = uint32(value.(float64))
			// TODO:
		case "visited_at":
			entity.VisitedAt = uint32(value.(float64))
		case "mark":
			entity.Mark = uint32(value.(float64))
		}
	}
	entity.Mutex.Unlock()
}

func (entity *Visit) to_json(w io.Writer) {
	entity.Mutex.RLock()
	json.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

func NewVisitsRepo() VisitsRepo {
	return VisitsRepo{
		Collection: make(map[uint32]*Visit),
		Mutex:      sync.RWMutex{}}
}

func (repo *VisitsRepo) Create(data *JsonData) {
	entity := &Visit{}
	entity.Update(data)
	repo.Add(entity)
}

func (repo *VisitsRepo) CreateFromJson(raw_data []byte) error {
	entity := &Visit{}
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

func (repo *VisitsRepo) Find(id uint32) (Entity, bool) {
	repo.Mutex.RLock()
	var entity, found = repo.Collection[id]
	repo.Mutex.RUnlock()
	return entity, found
}
