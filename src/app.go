package main

import (
	"archive/zip"
	"flag"
	"fmt"
	//"github.com/pquerna/ffjson/ffjson"
	"encoding/json"
	// "github.com/pkg/profile"
	"io/ioutil"
	"log"
	"strings"
	// "net/http"
	// "net/http/pprof"
	// "runtime/pprof"
	"sync"
)

var (
	ADDR      = flag.String("addr", ":80", "TCP address to listen to")
	DATA_PATH = flag.String("data", "/tmp/data/data.zip", "Initial data zipfile path")
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

func loadInitialData() {
	fmt.Println("DataLoading: starting")
	r, err := zip.OpenReader(DATA_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	entity_kinds := []string{"users", "locations", "visits"}
	for _, entity_kind := range entity_kinds {
		for _, f := range r.File {
			if !strings.Contains(f.Name, entity_kind) {
				continue
			}
			fmt.Println("DataLoading: loading", f.Name)

			rc, _ := f.Open()
			b, _ := ioutil.ReadAll(rc)
			data := make(JsonDataArray)
			json.Unmarshal(b, &data)
			json_objects := data[entity_kind]

			repo := entity_repo(len(entity_kind))
			for _, json_object := range json_objects {
				repo.Create(&json_object)
			}
		}
	}
}

func main() {
	flag.Parse()
	initVars()
	loadInitialData()

	// r := http.NewServeMux()
	// r.HandleFunc("/debug/pprof/", pprof.Index)
	// r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// go http.ListenAndServe(":8080", r)

	// defer profile.Start().Stop()
	startWebServer()
}
