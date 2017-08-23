package main

import (
	"fmt"
	//"github.com/pquerna/ffjson/ffjson"
	"encoding/json"
	"github.com/valyala/fasthttp"
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
	if lock {
		entity.Mutex.Unlock()
	}
}

func (entity *Location) toJson(w io.Writer) {
	entity.Mutex.RLock()
	json.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

func BirthDateToAge(BirthDate int) int {
	now := int(time.Now().Unix())
	age_ts := int64(now - BirthDate)
	age := int(time.Unix(age_ts, 0).Year() - 1970)
	return age
}

func AgeToBirthday(age int) int {
	// from datetime import datetime
	// from dateutil.relativedelta import relativedelta
	// import calendar

	birthday := time.Now().AddDate(-age, 0, 0)
	return int(birthday.Unix())
	// birthday = time.Date(now.Year(), time.November, 10, 23, 0, 0, 0, time.UTC)
	// now = datetime.now() - relativedelta(years = fromAge)

	// timestamp = calendar.timegm(now.timetuple())

	// 	now := int(time.Now().Unix())
	// 	age_ts := int64(now - BirthDate)
	// 	age := int(time.Unix(age_ts, 0).Year() - 1970)
	// 	return age
	// 	time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
}

// fromAge - учитывать только путешественников, у которых возраст (считается от текущего timestamp) строго больше этого параметра
// toAge - учитывать только путешественников, у которых возраст (считается от текущего timestamp) строго меньше этого параметра

// Небольшой пример проверки дат в этом запросе на python (fromAge - количество лет):

// Дальше проверяется birthdate < timestamp либо birthdate > timestamp соответственно.

func (entity *Location) checkVisit(visit *Visit, fromDate *int, toDate *int, fromAge *int, toAge *int, gender *string) bool {
	if fromDate != nil && visit.VisitedAt < *fromDate {
		return false
	}
	if toDate != nil && visit.VisitedAt > *toDate {
		return false
	}
	// if fromAge != nil || toAge != nil {
	// 	age := BirthDateToAge(visit.User.BirthDate)
	if fromAge != nil && visit.User.BirthDate >= AgeToBirthday(*fromAge) {
		return false
	}
	// if toAge != nil {
	// 	fmt.Println(*toAge, visit.User.BirthDate, AgeToBirthday(*toAge), BirthDateToAge(visit.User.BirthDate), BirthDateToAge(AgeToBirthday(*toAge)), (visit.User.BirthDate <= AgeToBirthday(*toAge)), visit.Mark)
	// }
	if toAge != nil && visit.User.BirthDate <= AgeToBirthday(*toAge) {
		return false
	}
	// }
	if gender != nil && visit.User.Gender != *gender {
		return false
	}
	return true
}

func (entity *Location) Visits(fromDate *int, toDate *int, fromAge *int, toAge *int, gender *string) []*Visit {
	visits_repo := LocationsVisits.findVisitsRepo(entity.Id)
	if visits_repo == nil {
		return nil
	}
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

func (entity *Location) WriteAvgsJson(w *fasthttp.RequestCtx, fromDate *int, toDate *int, fromAge *int, toAge *int, gender *string) {

	entity.Mutex.RLock()
	visits := entity.Visits(fromDate, toDate, fromAge, toAge, gender)
	entity.Mutex.RUnlock()

	if len(visits) == 0 {
		w.WriteString("{\"avg\": 0}")
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

		w.WriteString("{\"avg\": ")
		w.WriteString(avg_str)
		w.WriteString("}")
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
