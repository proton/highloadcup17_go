package main

import (
	"encoding/json"
	// "fmt"
	"io"
	"sort"
	"sync"
)

type User struct {
	Id          uint32          `json:"id"`
	Email       string          `json:"email"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	Gender      string          `json:"gender"`
	BirthDate   int32           `json:"birth_date"`
	Mutex       sync.RWMutex    `json:"-"`
	VisitIdsMap map[uint32]bool `json:"-"`
}

type UsersRepo struct {
	Collection map[uint32]*User
	Mutex      sync.RWMutex
}

func (entity *User) Update(data *JsonData, lock bool) bool {
	if lock {
		entity.Mutex.Lock()
		defer entity.Mutex.Unlock()
	}
	denormolize_in_visits := false
	for key, value := range *data {
		if value == nil {
			return false
		}
		switch key {
		case "id":
			entity.Id = uint32(value.(float64))
		case "email":
			entity.Email = value.(string)
		case "first_name":
			entity.FirstName = value.(string)
		case "last_name":
			entity.LastName = value.(string)
		case "gender":
			gender := value.(string)
			if !validate_gender(gender) {
				return false
			}
			entity.Gender = gender
			denormolize_in_visits = true
		case "birth_date":
			entity.BirthDate = int32(value.(float64))
			denormolize_in_visits = true
		}
	}
	if denormolize_in_visits {
		visits := entity.Visits(nil, nil, nil, nil)
		for _, visit := range visits {
			visit.UserGender = entity.Gender
			visit.UserBirthDate = entity.BirthDate
			visit.Mutex.RUnlock()
		}
	}
	return true
}

func (entity *User) to_json(w io.Writer) {
	entity.Mutex.RLock()
	json.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

func (entity *User) VisitIds() []uint32 {
	ids := make([]uint32, len(entity.VisitIdsMap))

	i := 0
	for id := range entity.VisitIdsMap {
		ids[i] = id
		i++
	}
	return ids
}

func (entity *User) Visits(fromDate *uint32, toDate *uint32, country *string, toDistance *uint32) []*Visit {
	visits, _ := Visits.FindAll(entity.VisitIds())
	filteredVisits := make([]*Visit, 0, len(visits))
	for _, visit := range visits {
		visit.Mutex.RLock()
		if visit.UserId != entity.Id {
			visit.Mutex.RUnlock()
			continue
		}
		if fromDate != nil && visit.VisitedAt < *fromDate {
			visit.Mutex.RUnlock()
			continue
		}
		if toDate != nil && visit.VisitedAt > *toDate {
			visit.Mutex.RUnlock()
			continue
		}
		if country != nil && visit.LocationCountry != *country {
			visit.Mutex.RUnlock()
			continue
		}
		if toDistance != nil && visit.LocationDistance >= *toDistance {
			visit.Mutex.RUnlock()
			continue
		}
		filteredVisits = append(filteredVisits, visit)
	}
	return filteredVisits
}

func (entity *User) WriteVisitsJson(w io.Writer, fromDate *uint32, toDate *uint32, country *string, toDistance *uint32) {
	entity.Mutex.RLock()

	visits := entity.Visits(fromDate, toDate, country, toDistance)

	sort.Slice(visits, func(i, j int) bool { return visits[i].VisitedAt < visits[j].VisitedAt })

	w.Write([]byte("{\"visits\": ["))
	first := true
	for _, visit := range visits {
		if first == false {
			w.Write([]byte(","))
		}
		json.NewEncoder(w).Encode(visit.ToView())
		visit.Mutex.RUnlock()
		first = false
	}
	w.Write([]byte("]}"))
	entity.Mutex.RUnlock()
}

func NewUsersRepo() UsersRepo {
	return UsersRepo{
		Collection: make(map[uint32]*User),
		Mutex:      sync.RWMutex{}}
}

func (repo *UsersRepo) InitEntity() *User {
	entity := User{
		VisitIdsMap: make(map[uint32]bool)}
	return &entity
}

func (repo *UsersRepo) Create(data *JsonData) bool {
	entity := repo.InitEntity()
	ok := entity.Update(data, false)
	if !ok {
		return false
	}
	repo.Add(entity)
	return true
}

// func (repo *UsersRepo) CreateFromJson(raw_data []byte) error {
// 	entity := repo.InitEntity()
// 	err := json.Unmarshal(raw_data, entity)
// 	if err == nil {
// 		repo.Add(entity)
// 	}
// 	return err
// }

func (repo *UsersRepo) Add(entity *User) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *UsersRepo) Find(id uint32, lock bool) (*User, bool) {
	if lock {
		repo.Mutex.RLock()
	}
	var entity, found = repo.Collection[id]
	if lock {
		repo.Mutex.RUnlock()
	}
	return entity, found
}

func (repo *UsersRepo) FindEntity(id uint32, lock bool) (Entity, bool) {
	return repo.Find(id, lock)
}
