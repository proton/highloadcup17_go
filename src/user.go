package main

import (
	// "fmt"
	"github.com/pquerna/ffjson/ffjson"
	"io"
	"sort"
	"strconv"
	"sync"
)

type User struct {
	Id        int          `json:"id"`
	Email     string       `json:"email"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Gender    string       `json:"gender"`
	BirthDate int          `json:"birth_date"`
	Mutex     sync.RWMutex `json:"-"`
}

type UsersRepo struct {
	Collection map[int]*User
	Mutex      sync.RWMutex
}

func (entity *User) Update(data *JsonData, lock bool) {
	if lock {
		entity.Mutex.Lock()
		defer entity.Mutex.Unlock()
	}
	for key, value := range *data {
		switch key {
		case "id":
			entity.Id = int(value.(float64))
		case "email":
			entity.Email = value.(string)
		case "first_name":
			entity.FirstName = value.(string)
		case "last_name":
			entity.LastName = value.(string)
		case "gender":
			entity.Gender = value.(string)
		case "birth_date":
			entity.BirthDate = int(value.(float64))
		}
	}
}

func (entity *User) to_json(w io.Writer) {
	entity.Mutex.RLock()
	ffjson.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

func (entity *User) checkVisit(visit *Visit, fromDate *int, toDate *int, country *string, toDistance *int) bool {
	if visit.UserId != entity.Id {
		return false
	}
	if fromDate != nil && visit.VisitedAt < *fromDate {
		return false
	}
	if toDate != nil && visit.VisitedAt > *toDate {
		return false
	}
	if country != nil && visit.Location.Country != *country {
		return false
	}
	if toDistance != nil && visit.Location.Distance >= *toDistance {
		return false
	}
	return true
}

func (entity *User) Visits(fromDate *int, toDate *int, country *string, toDistance *int) []*Visit {
	visits_repo := UsersVisits.findVisitsRepo(entity.Id)
	if visits_repo == nil {
		return nil
	}
	visits_repo.Mutex.RLock()
	visits := make([]*Visit, 0, len(visits_repo.Collection))
	for _, visit := range visits_repo.Collection {
		visit.Mutex.RLock()
		if !entity.checkVisit(visit, fromDate, toDate, country, toDistance) {
			continue
		}
		visits = append(visits, visit)
		visit.Mutex.RUnlock()
	}
	visits_repo.Mutex.RUnlock()
	return visits
}

func (entity *User) WriteVisitsJson(w io.Writer, fromDate *int, toDate *int, country *string, toDistance *int) {

	entity.Mutex.RLock()
	visits := entity.Visits(fromDate, toDate, country, toDistance)
	entity.Mutex.RUnlock()

	for _, visit := range visits {
		visit.Mutex.RLock()
	}

	sort.Slice(visits, func(i, j int) bool { return visits[i].VisitedAt < visits[j].VisitedAt })

	w.Write([]byte("{\"visits\": ["))
	first := true
	for _, visit := range visits {
		if first == false {
			w.Write([]byte(","))
		}
		ffjson.NewEncoder(w).Encode(visit.ToView())
		first = false
	}
	w.Write([]byte("]}"))

	for _, visit := range visits {
		visit.Mutex.RUnlock()
	}
}

func (repo *UsersRepo) InitEntity() *User {
	return &User{}
}

func (repo *UsersRepo) Create(data *JsonData) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
}

func (repo *UsersRepo) Add(entity *User) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *UsersRepo) Find(id int) (*User, bool) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	entity, ok := repo.Collection[id]
	return entity, ok
}

func (repo *UsersRepo) FindEntity(id int) (Entity, bool) {
	return repo.Find(id)
}

func find_user(entity_id_str *string) (*User, bool) {
	entity_id_int, error := strconv.Atoi(*entity_id_str)
	if error == nil {
		entity_id := int(entity_id_int)
		return Users.Find(entity_id)
	}
	return nil, false
}
