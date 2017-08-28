package main

import (
	"flag"
	// "github.com/pkg/profile"
	"sync"
	"time"
)

func main() {
	flag.Parse()
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

	Users = UsersRepo{
		Collection:    make([]*User, USERS_REPO_COLLECTION_SIZE),
		MapCollection: make(map[uint32]*User),
		MapMutex:      sync.RWMutex{}}
	Locations = LocationsRepo{
		Collection:    make([]*Location, LOCATIONS_REPO_COLLECTION_SIZE),
		MapCollection: make(map[uint32]*Location),
		MapMutex:      sync.RWMutex{}}
	Visits = VisitsRepo{
		Collection:    make([]*Visit, VISITS_REPO_COLLECTION_SIZE),
		MapCollection: make(map[uint32]*Visit),
		MapMutex:      sync.RWMutex{}}
	UsersVisits = EntityVisitsRepo{
		Collection: make(map[uint32]*VisitsMRepo),
		Mutex:      sync.RWMutex{}}
	LocationsVisits = EntityVisitsRepo{
		Collection: make(map[uint32]*VisitsMRepo),
		Mutex:      sync.RWMutex{}}

	InitialTime time.Time
)

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
