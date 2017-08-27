package main

import (
	"flag"
	// "github.com/pkg/profile"
	"sync"
	"time"
)

var (
	ADDR         = flag.String("addr", ":80", "TCP address to listen to")
	DATAZIP_PATH = flag.String("zip", "/tmp/data/data.zip", "Zipfile path")
	DATA_DIR     = flag.String("data", "/", "Directory with extacted jsons")
)

var (
	Users           UsersRepo
	Locations       LocationsRepo
	Visits          VisitsRepo
	UsersVisits     EntityVisitsRepo
	LocationsVisits EntityVisitsRepo
	InitialTime     time.Time
)

func initVars() {
	Users = UsersRepo{
		Collection: make(map[uint32]*User),
		Mutex:      sync.RWMutex{}}
	Locations = LocationsRepo{
		Collection: make(map[uint32]*Location),
		Mutex:      sync.RWMutex{}}
	Visits = VisitsRepo{
		Collection: make(map[uint32]*Visit),
		Mutex:      sync.RWMutex{}}
	UsersVisits = EntityVisitsRepo{
		Collection: make(map[uint32]*VisitsRepo),
		Mutex:      sync.RWMutex{}}
	LocationsVisits = EntityVisitsRepo{
		Collection: make(map[uint32]*VisitsRepo),
		Mutex:      sync.RWMutex{}}
}

func entity_repo(entity_kind_len int) EntityRepo {
	switch entity_kind_len {
	case 5: //"users":
		return &Users
	case 9: //"locations":
		return &Locations
	case 6: //"visits":
		return &Visits
	}
	return nil
}

func main() {
	flag.Parse()
	initVars()
	loadInitialData()
	startWebServer()
}

// func main() {
// 	flag.Parse()
// 	initVars()
// 	profile.Start(profile.MemProfile)
// 	// profile.Start(profile.CPUProfile)
// 	// profile.Start(profile.MutexProfile)
// 	// profile.Start(profile.BlockProfile)
// 	loadInitialData()
// 	// defer profile.Start().Stop()
// 	startWebServer()
// }
