package main

import (
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"io"
	"strconv"
	"sync"
	"time"
)

type Location struct {
	Id       int          `json:"id"`
	Place    string       `json:"place"`
	Country  string       `json:"country"`
	City     string       `json:"city"`
	Distance int          `json:"distance"`
	Mutex    sync.RWMutex `json:"-"`
}

type LocationsRepo struct {
	Collection map[int]*Location
	Mutex      sync.RWMutex
}

func (entity *Location) Update(data *JsonData, lock bool) {
	if lock {
		entity.Mutex.Lock()
		defer entity.Mutex.Unlock()
	}
	for key, value := range *data {
		switch key {
		case "id":
			entity.Id = int(value.(float64))
		case "place":
			entity.Place = value.(string)
		case "country":
			entity.Country = value.(string)
		case "city":
			entity.City = value.(string)
		case "distance":
			entity.Distance = int(value.(float64))
		}
	}
}

func (entity *Location) to_json(w io.Writer) {
	entity.Mutex.RLock()
	ffjson.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

// func (entity *Location) VisitIds() []int {
// 	ids := make([]int, len(entity.VisitIdsMap))

// 	i := 0
// 	for id := range entity.VisitIdsMap {
// 		ids[i] = id
// 		i++
// 	}
// 	return ids
// }

func BirthDateToAge(BirthDate int) int {
	now := int(time.Now().Unix())
	age_ts := int64(now - BirthDate)
	age := int(time.Unix(age_ts, 0).Year() - 1970)
	return age
}

func (entity *Location) checkVisit(visit *Visit, fromDate *int, toDate *int, fromAge *int, toAge *int, gender *string) bool {
	if fromDate != nil && visit.VisitedAt < *fromDate {
		return false
	}
	if toDate != nil && visit.VisitedAt > *toDate {
		return false
	}
	if fromAge != nil || toAge != nil {
		age := BirthDateToAge(visit.User.BirthDate)
		if fromAge != nil && age <= *fromAge {
			return false
		}
		if toAge != nil && age >= *toAge {
			return false
		}
	}
	if gender != nil && visit.User.Gender != *gender {
		return false
	}
	return true
}

func (entity *Location) Visits(fromDate *int, toDate *int, fromAge *int, toAge *int, gender *string) []*Visit {
	visits_repo := LocationsVisits.findVisitsRepo(entity.Id)
	visits_repo.Mutex.RLock()
	filteredVisits := make([]*Visit, 0, len(visits_repo.Collection))
	for _, visit := range visits_repo.Collection {
		visit.Mutex.RLock()
		if !entity.checkVisit(visit, fromDate, toDate, fromAge, toAge, gender) {
			continue
		}
		filteredVisits = append(filteredVisits, visit)
		visit.Mutex.RUnlock()
	}
	visits_repo.Mutex.RUnlock()
	return filteredVisits
}

func (entity *Location) WriteAvgsJson(w io.Writer, fromDate *int, toDate *int, fromAge *int, toAge *int, gender *string) {

	entity.Mutex.RLock()
	visits := entity.Visits(fromDate, toDate, fromAge, toAge, gender)
	entity.Mutex.RUnlock()

	if len(visits) == 0 {
		w.Write([]byte("{\"avg\": 0}"))
	} else {
		marks_count := 0
		marks_sum := 0

		for _, visit := range visits {
			visit.Mutex.RLock()
			marks_sum += visit.Mark
			marks_count += 1
			visit.Mutex.RUnlock()
		}

		avg := float64(marks_sum) / float64(marks_count)
		avg_str := fmt.Sprintf("%.5f", avg)

		w.Write([]byte("{\"avg\": "))
		w.Write([]byte(avg_str))
		w.Write([]byte("}"))
	}
}

func (repo *LocationsRepo) InitEntity() *Location {
	return &Location{}
}

func (repo *LocationsRepo) Create(data *JsonData) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
}

func (repo *LocationsRepo) Add(entity *Location) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *LocationsRepo) Find(id int) (*Location, bool) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	entity, ok := repo.Collection[id]
	return entity, ok
}

func (repo *LocationsRepo) FindEntity(id int) (Entity, bool) {
	return repo.Find(id)
}

func find_location(entity_id_str *string) (*Location, bool) {
	entity_id_int, error := strconv.Atoi(*entity_id_str)
	if error == nil {
		entity_id := int(entity_id_int)
		return Locations.Find(entity_id)
	}
	return nil, false
}
