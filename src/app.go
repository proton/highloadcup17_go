package main

import (
	"flag"
	// "github.com/pquerna/ffjson/ffjson"
	// "encoding/json"
	// "github.com/pkg/profile"
	// "net/http"
	// "net/http/pprof"
	// "runtime/pprof"
	"sync"
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

	// r := http.NewServeMux()
	// r.HandleFunc("/debug/pprof/", pprof.Index)
	// r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// go http.ListenAndServe(":8080", r)

	loadInitialData()

	// defer func() {
	// 	if fd, err := os.Create(`pprof.mem`); err == nil {
	// 		pprof.WriteHeapProfile(fd)
	// 		fd.Close()
	// 	}
	// }()

	// defer profile.Start().Stop()
	startWebServer()
}
