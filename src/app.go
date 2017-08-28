package main

import (
	"flag"
	// "github.com/pkg/profile"
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

	Users           = makeUsersRepo(USERS_REPO_COLLECTION_SIZE)
	Locations       = makeLocationsRepo(LOCATIONS_REPO_COLLECTION_SIZE)
	Visits          = makeVisitsRepo(VISITS_REPO_COLLECTION_SIZE)
	UsersVisits     = makeEntityVisitsRepo(USERS_REPO_COLLECTION_SIZE)
	LocationsVisits = makeEntityVisitsRepo(LOCATIONS_REPO_COLLECTION_SIZE)

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
