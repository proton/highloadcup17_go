package main

import (
	"github.com/buger/jsonparser"
	"github.com/pquerna/ffjson/ffjson"
	// "fmt"
	"io"
	"sync"
)

type Visit struct {
	Id         uint32       `json:"id"`
	LocationId uint32       `json:"location"`
	UserId     uint32       `json:"user"`
	VisitedAt  uint32       `json:"visited_at"`
	Mark       uint8        `json:"mark"`
	Mutex      sync.RWMutex `json:"-"`
	Location   *Location    `json:"-"`
	User       *User        `json:"-"`
	Json       []byte       `json:"-"`
}

type VisitView struct {
	Mark      uint8  `json:"mark"`
	VisitedAt uint32 `json:"visited_at"`
	Place     string `json:"place"`
}

var (
	VISIT_JSON_PATHS = [][]string{
		[]string{"id"},
		[]string{"location"},
		[]string{"user"},
		[]string{"visited_at"},
		[]string{"mark"},
	}
)

func (entity *Visit) UpdateFromJSON(data []byte, lock bool) {
	sync_user := false
	sync_location := false
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
			if v, er := jsonparser.ParseInt(value); er == nil {
				entity.LocationId = uint32(v)
				sync_location = true
				location, _ := Locations.Find(entity.LocationId)
				entity.Location = location
			}
		case 2:
			if v, er := jsonparser.ParseInt(value); er == nil {
				entity.UserId = uint32(v)
				sync_user = true
				user, _ := Users.Find(entity.UserId)
				entity.User = user
			}
		case 3:
			if v, er := jsonparser.ParseInt(value); er == nil {
				entity.VisitedAt = uint32(v)
			}
		case 4:
			if v, er := jsonparser.ParseInt(value); er == nil {
				entity.Mark = uint8(v)
			}
		}
	}, VISIT_JSON_PATHS...)

	// entity.cacheJSON() // Too much memory :(
	if lock {
		entity.Mutex.Unlock()
	}

	if sync_location {
		LocationsVisits.addVisit(entity.LocationId, entity)
	}
	if sync_user {
		UsersVisits.addVisit(entity.UserId, entity)
	}
}

func (entity *Visit) cacheJSON() {
	b, _ := ffjson.Marshal(entity)
	entity.Json = b
}

func (entity *Visit) writeJSON(w io.Writer) {
	entity.Mutex.RLock()
	defer entity.Mutex.RUnlock()
	ffjson.NewEncoder(w).Encode(entity)
}

func (visit *Visit) ToView() *VisitView {
	return &VisitView{
		Mark:      visit.Mark,
		VisitedAt: visit.VisitedAt,
		Place:     visit.Location.Place}
}
