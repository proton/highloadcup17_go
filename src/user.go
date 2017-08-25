package main

import (
	// "fmt"
	"github.com/pquerna/ffjson/ffjson"
	// "encoding/json"
	"github.com/valyala/fasthttp"
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
	Json      []byte       `json:"-"`
}

type UsersRepo struct {
	Collection map[int]*User
	Mutex      sync.RWMutex
}

func (entity *User) Update(data *JsonData, lock bool) {
	if lock {
		entity.Mutex.Lock()
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
	entity.cacheJSON()
	if lock {
		entity.Mutex.Unlock()
	}
}

func (entity *User) cacheJSON() {
	b, _ := ffjson.Marshal(entity)
	entity.Json = b
}

func (entity *User) writeJSON(w io.Writer) {
	entity.Mutex.RLock()
	w.Write(entity.Json)
	entity.Mutex.RUnlock()
}

func (entity *User) checkVisit(visit *Visit, fromDate *int, toDate *int, country *string, toDistance *int) bool {
	visit.Mutex.RLock()
	defer visit.Mutex.RUnlock()
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
		if !entity.checkVisit(visit, fromDate, toDate, country, toDistance) {
			continue
		}
		visits = append(visits, visit)
	}
	visits_repo.Mutex.RUnlock()
	return visits
}

func (entity *User) WriteVisitsJson(w *fasthttp.RequestCtx, fromDate *int, toDate *int, country *string, toDistance *int) {

	entity.Mutex.RLock()
	visits := entity.Visits(fromDate, toDate, country, toDistance)
	entity.Mutex.RUnlock()

	for _, visit := range visits {
		visit.Mutex.RLock()
	}

	sort.Slice(visits, func(i, j int) bool { return visits[i].VisitedAt < visits[j].VisitedAt })

	w.WriteString("{\"visits\": [")
	first := true
	for _, visit := range visits {
		if first == false {
			w.WriteString(",")
		}
		ffjson.NewEncoder(w).Encode(visit.ToView())
		first = false
	}
	w.WriteString("]}")

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
