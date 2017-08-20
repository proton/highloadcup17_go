package main

import (
	"encoding/json"
	"io"
	"sync"
)

type Location struct {
	Id       uint32       `json:"id"`
	Place    string       `json:"place"`
	Country  string       `json:"country"`
	City     string       `json:"city"`
	Distance uint32       `json:"distance"`
	Mutex    sync.RWMutex `json:"-"`
	// VisitIds []uint32
}

type LocationsRepo struct {
	Collection map[uint32]*Location
	Mutex      sync.RWMutex
}

func (entity *Location) Update(data *JsonData) {
	entity.Mutex.Lock()
	for key, value := range *data {
		switch key {
		case "id":
			entity.Id = uint32(value.(float64))
		case "place":
			entity.Place = value.(string)
		case "country":
			entity.Country = value.(string)
		case "city":
			entity.City = value.(string)
		case "distance":
			entity.Distance = uint32(value.(float64))
		}
	}
	entity.Mutex.Unlock()
}

func (entity *Location) to_json(w io.Writer) {
	entity.Mutex.RLock()
	json.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

func NewLocationsRepo() LocationsRepo {
	return LocationsRepo{
		Collection: make(map[uint32]*Location),
		Mutex:      sync.RWMutex{}}
}

func (repo *LocationsRepo) Create(data *JsonData) {
	entity := &Location{}
	entity.Update(data)
	repo.Add(entity)
}

func (repo *LocationsRepo) CreateFromJson(raw_data []byte) error {
	entity := &Location{}
	err := json.Unmarshal(raw_data, entity)
	if err == nil {
		repo.Add(entity)
	}
	return err
}

func (repo *LocationsRepo) Add(entity *Location) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *LocationsRepo) Find(id uint32) (Entity, bool) {
	repo.Mutex.RLock()
	var entity, found = repo.Collection[id]
	repo.Mutex.RUnlock()
	return entity, found
}
