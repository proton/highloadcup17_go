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
	"unsafe"
)

var (
	ADDR     = flag.String("addr", ":80", "TCP address to listen to")
	DATA_DIR = flag.String("data", "/tmp/data/", "Directory with zipfile path")
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
		Collection: make(map[int]*User),
		Mutex:      sync.RWMutex{}}
	Locations = LocationsRepo{
		Collection: make(map[int]*Location),
		Mutex:      sync.RWMutex{}}
	Visits = VisitsRepo{
		Collection: make(map[int]*Visit),
		Mutex:      sync.RWMutex{}}
	UsersVisits = EntityVisitsRepo{
		Collection: make(map[int]*VisitsRepo),
		Mutex:      sync.RWMutex{}}
	LocationsVisits = EntityVisitsRepo{
		Collection: make(map[int]*VisitsRepo),
		Mutex:      sync.RWMutex{}}
}

func bstring(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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

	// defer profile.Start().Stop()
	startWebServer()
}
