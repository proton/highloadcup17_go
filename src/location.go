package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"io"
	"strconv"
	"sync"
)

type Location struct {
	Id       uint32       `json:"id"`
	Place    string       `json:"place"`
	Country  string       `json:"country"`
	City     string       `json:"city"`
	Distance uint32       `json:"distance"`
	Mutex    sync.RWMutex `json:"-"`
	Json     []byte       `json:"-"`
}

type LocationsRepo struct {
	Collection map[uint32]*Location
	Mutex      sync.RWMutex
}

var (
	LOCATION_JSON_PATHS = [][]string{
		[]string{"id"},
		[]string{"place"},
		[]string{"country"},
		[]string{"city"},
		[]string{"distance"},
	}
)

func (entity *Location) UpdateFromJSON(data []byte, lock bool) {
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
				entity.Place = v
			}
		case 2:
			if v, er := jsonparser.ParseString(value); er == nil {
				entity.Country = v
			}
		case 3:
			if v, er := jsonparser.ParseString(value); er == nil {
				entity.City = v
			}
		case 4:
			if v, er := jsonparser.ParseInt(value); er == nil {
				entity.Distance = uint32(v)
			}
		}
	}, LOCATION_JSON_PATHS...)

	entity.cacheJSON()
	if lock {
		entity.Mutex.Unlock()
	}
}

func (entity *Location) cacheJSON() {
	b, _ := ffjson.Marshal(entity)
	entity.Json = b
}

func (entity *Location) writeJSON(w io.Writer) {
	entity.Mutex.RLock()
	defer entity.Mutex.RUnlock()
	w.Write(entity.Json)
}

func (entity *Location) checkVisit(visit *Visit, fromDate *uint32, toDate *uint32, fromAgeBirthday *int32, toAgeBirthday *int32, gender *string) bool {
	visit.Mutex.RLock()
	defer visit.Mutex.RUnlock()
	if visit.LocationId != entity.Id {
		return false
	}
	if fromDate != nil && visit.VisitedAt < *fromDate {
		return false
	}
	if toDate != nil && visit.VisitedAt > *toDate {
		return false
	}
	if fromAgeBirthday != nil && visit.User.BirthDate >= *fromAgeBirthday {
		return false
	}
	if toAgeBirthday != nil && visit.User.BirthDate <= *toAgeBirthday {
		return false
	}
	if gender != nil && visit.User.Gender != *gender {
		return false
	}
	return true
}

func (entity *Location) WriteAvgsJson(w *fasthttp.RequestCtx, fromDate *uint32, toDate *uint32, fromAge *uint32, toAge *uint32, gender *string) {
	marks_count := 0
	marks_sum := uint32(0)

	visits_repo := LocationsVisits.findVisitsRepo(entity.Id)
	if visits_repo != nil {
		fromAgeBirthday := AgeToBirthday(fromAge)
		toAgeBirthday := AgeToBirthday(toAge)

		visits_repo.Mutex.RLock()
		for _, visit := range visits_repo.Collection {
			if !entity.checkVisit(visit, fromDate, toDate, fromAgeBirthday, toAgeBirthday, gender) {
				continue
			}
			marks_count += 1
			visit.Mutex.RLock()
			marks_sum += uint32(visit.Mark)
			visit.Mutex.RUnlock()
		}
		visits_repo.Mutex.RUnlock()
	}

	if marks_count == 0 {
		w.WriteString("{\"avg\": 0}")
	} else {
		avg := float64(marks_sum) / float64(marks_count)
		avg_str := fmt.Sprintf("%.5f", avg)

		w.WriteString("{\"avg\": ")
		w.WriteString(avg_str)
		w.WriteString("}")
	}
}

func (repo *LocationsRepo) InitEntity() *Location {
	return &Location{}
}

func (repo *LocationsRepo) CreateFromJSON(data []byte) {
	entity := repo.InitEntity()
	entity.UpdateFromJSON(data, false)
	repo.Add(entity)
}

func (repo *LocationsRepo) Add(entity *Location) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *LocationsRepo) Find(id uint32) (*Location, bool) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	entity, ok := repo.Collection[id]
	return entity, ok
}

func (repo *LocationsRepo) FindEntity(id uint32) (Entity, bool) {
	return repo.Find(id)
}

func find_location(entity_id_str []byte) (*Location, bool) {
	entity_id_int, error := strconv.Atoi(bstring(entity_id_str))
	if error == nil {
		entity_id := uint32(entity_id_int)
		return Locations.Find(entity_id)
	}
	return nil, false
}
