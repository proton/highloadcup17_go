package main

import (
	// "fmt"
	"github.com/buger/jsonparser"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"io"
	"sort"
	"sync"
)

type User struct {
	Id        uint32       `json:"id"`
	Email     string       `json:"email"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Gender    string       `json:"gender"`
	BirthDate int32        `json:"birth_date"`
	Mutex     sync.RWMutex `json:"-"`
	Json      []byte       `json:"-"`
}

var (
	USER_JSON_PATHS = [][]string{
		[]string{"id"},
		[]string{"email"},
		[]string{"first_name"},
		[]string{"last_name"},
		[]string{"gender"},
		[]string{"birth_date"},
	}
)

func (entity *User) UpdateFromJSON(data []byte, lock bool) {
	if lock {
		entity.Mutex.Lock()
	}

	jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		switch idx {
		case 0:
			if v, er := jsonparser.ParseInt(value); er == nil {
				entity.Id = uint32(v)
			}
		case 1:
			if v, er := jsonparser.ParseString(value); er == nil {
				entity.Email = v
			}
		case 2:
			if v, er := jsonparser.ParseString(value); er == nil {
				entity.FirstName = v
			}
		case 3:
			if v, er := jsonparser.ParseString(value); er == nil {
				entity.LastName = v
			}
		case 4:
			if v, er := jsonparser.ParseString(value); er == nil {
				entity.Gender = v
			}
		case 5:
			if v, er := jsonparser.ParseInt(value); er == nil {
				entity.BirthDate = int32(v)
			}
		}
	}, USER_JSON_PATHS...)

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
	defer entity.Mutex.RUnlock()
	w.Write(entity.Json)
}

func (entity *User) checkVisit(visit *Visit, fromDate *uint32, toDate *uint32, country *string, toDistance *uint32) bool {
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

func (entity *User) Visits(fromDate *uint32, toDate *uint32, country *string, toDistance *uint32) []*Visit {
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

func (entity *User) WriteVisitsJson(w *fasthttp.RequestCtx, fromDate *uint32, toDate *uint32, country *string, toDistance *uint32) {

	entity.Mutex.RLock()
	visits := entity.Visits(fromDate, toDate, country, toDistance)
	entity.Mutex.RUnlock()

	for _, visit := range visits {
		visit.Mutex.RLock()
	}

	sort.Slice(visits, func(i, j int) bool { return visits[i].VisitedAt < visits[j].VisitedAt })

	w.WriteString("{\"visits\": [")
	first := true
	enc := ffjson.NewEncoder(w)
	for _, visit := range visits {
		if first == false {
			w.WriteString(",")
		}
		enc.Encode(visit.ToView())
		first = false
	}
	w.WriteString("]}")

	for _, visit := range visits {
		visit.Mutex.RUnlock()
	}
}
