package main

import (
	"flag"
	// "github.com/pkg/profile"
	"sync"
	"time"
)

func main() {
	flag.Parse()
	initVars()
	loadInitialData()
	startWebServer()
}

// func main() {
// 	flag.Parse()
// 	initVars()
// 	// defer profile.Start(profile.CPUProfile).Stop()
// 	// defer profile.Start(profile.MutexProfile).Stop()
// 	// defer profile.Start(profile.BlockProfile).Stop()
// 	loadInitialData()
// 	defer profile.Start(profile.MemProfile).Stop()
// 	startWebServer()
// }

var (
	ADDR         = flag.String("addr", ":80", "TCP address to listen to")
	DATAZIP_PATH = flag.String("zip", "/tmp/data/data.zip", "Zipfile path")
	DATA_DIR     = flag.String("data", "/", "Directory with extacted jsons")
	OPTIONS_PATH = flag.String("options", "/tmp/data/options.txt", "options file path")
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
		Collection: make([]*Visit, VISITS_REPO_COLLECTION_SIZE),
		// Mutex:         sync.RWMutex{},
		MapCollection: make(map[uint32]*Visit),
		MapMutex:      sync.RWMutex{}}
	UsersVisits = EntityVisitsRepo{
		Collection: make(map[uint32]*VisitsMRepo),
		Mutex:      sync.RWMutex{}}
	LocationsVisits = EntityVisitsRepo{
		Collection: make(map[uint32]*VisitsMRepo),
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
