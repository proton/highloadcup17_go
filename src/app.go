package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	// "github.com/hashicorp/go-memdb"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"strings"
	// "time"
	// "github.com/pkg/profile"
	// "runtime/debug"
	"strconv"
)

var (
	addr = flag.String("addr", ":9000", "TCP address to listen to")
)

var (
	Users     UsersRepo
	Locations LocationsRepo
	Visits    VisitsRepo
)

func entity_repo(entity_kind string) EntityRepo {
	switch entity_kind {
	case "users":
		return &Users
	case "locations":
		return &Locations
	case "visits":
		return &Visits
	}
	return nil
}

func loadInitialData() {
	fmt.Println("DataLoading: starting")
	r, err := zip.OpenReader("/tmp/data/data.zip")
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
			ffjson.Unmarshal(b, &data)
			json_objects := data[entity_kind]

			repo := entity_repo(entity_kind)
			for _, json_object := range json_objects {
				repo.Create(&json_object)
			}
		}
	}
}

func startWebServer() {
	fmt.Println("Webserver: starting")
	h := requestHandler
	// h = fasthttp.CompressHandler(h)

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func main() {
	flag.Parse()

	Users = NewUsersRepo()
	Locations = NewLocationsRepo()
	Visits = NewVisitsRepo()

	loadInitialData()
	// defer profile.Start().Stop()
	startWebServer()
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	// ctx.SetContentType("text/plain; charset=utf8")
	defer func() {
		if r := recover(); r != nil {
			// fmt.Fprintf(ctx, "\n\nWEB SERVER ERROR: %s\n%s\n", r, debug.Stack())
			render400(ctx)
		}
	}()

	path := strings.Split(string(ctx.Path()), "/")
	http_method_is_get := string(ctx.Method()) == "GET"
	path_len := len(path)
	entity_kind := path[1]
	repo := entity_repo(entity_kind)

	if path_len == 3 {
		if http_method_is_get {
			var entity, ok = find_entity(repo, &path[2], false)
			if ok == true {
				renderEntity(ctx, entity)
				return
			}
		} else if path[2] == "new" {
			processEntityCreate(ctx, repo)
			return
		} else {
			var entity, ok = find_entity(repo, &path[2], true)
			if ok == true {
				processEntityUpdate(ctx, entity)
				return
			}
		}
	} else if path_len == 4 && http_method_is_get {
		if entity_kind == "users" && path[3] == "visits" {
			var user, ok = find_user(&path[2], false)
			if ok == true {
				processUserVisits(ctx, user)
				return
			}
		} else if entity_kind == "locations" && path[3] == "avg" {
			var location, ok = find_location(&path[2], false)
			if ok == true {
				processLocationAvgs(ctx, location)
				return
			}
		}
	}
	render404(ctx)
}

func find_entity(repo EntityRepo, entity_id_str *string, lock bool) (Entity, bool) {
	entity_id_int, error := strconv.Atoi(*entity_id_str)
	if error == nil {
		entity_id := uint32(entity_id_int)
		return repo.FindEntity(entity_id, lock)
	}
	return nil, false
}

func find_user(entity_id_str *string, lock bool) (*User, bool) {
	entity_id_int, error := strconv.Atoi(*entity_id_str)
	if error == nil {
		entity_id := uint32(entity_id_int)
		return Users.Find(entity_id, lock)
	}
	return nil, false
}

func find_location(entity_id_str *string, lock bool) (*Location, bool) {
	entity_id_int, error := strconv.Atoi(*entity_id_str)
	if error == nil {
		entity_id := uint32(entity_id_int)
		return Locations.Find(entity_id, lock)
	}
	return nil, false
}

func extractUintParam(ctx *fasthttp.RequestCtx, key string) (*uint32, bool) {
	param := ctx.QueryArgs().Peek(key)
	if param == nil {
		return nil, true
	}
	param_int, err := strconv.Atoi(string(param))
	if err != nil {
		return nil, false
	}
	param_uint := uint32(param_int)
	return &param_uint, true
}

func extractStringParam(ctx *fasthttp.RequestCtx, key string) (*string, bool) {
	param := ctx.QueryArgs().Peek(key)
	if param == nil {
		return nil, true
	}
	param_string := string(param)
	return &param_string, true
}

func processLocationAvgs(ctx *fasthttp.RequestCtx, location *Location) {
	fromDate, ok := extractUintParam(ctx, "fromDate")
	if ok {
		toDate, ok := extractUintParam(ctx, "toDate")
		if ok {
			fromAge, ok := extractUintParam(ctx, "fromAge")
			if ok {
				toAge, ok := extractUintParam(ctx, "toAge")
				if ok {
					gender, ok := extractStringParam(ctx, "gender")
					if ok && (gender == nil || validate_gender(*gender)) {
						location.WriteAvgsJson(ctx, fromDate, toDate, fromAge, toAge, gender)
						return
					}
				}
			}
		}
	}
	render400(ctx)
}

func processUserVisits(ctx *fasthttp.RequestCtx, user *User) {
	fromDate, ok := extractUintParam(ctx, "fromDate")
	if ok {
		toDate, ok := extractUintParam(ctx, "toDate")
		if ok {
			country, ok := extractStringParam(ctx, "country")
			if ok {
				toDistance, ok := extractUintParam(ctx, "toDistance")
				if ok {
					user.WriteVisitsJson(ctx, fromDate, toDate, country, toDistance)
					return
				}
			}
		}
	}
	render400(ctx)
}

type JsonData map[string]interface{}
type JsonDataArray map[string][]JsonData

func loadJSON(ctx *fasthttp.RequestCtx) *JsonData {
	var data JsonData
	body := ctx.PostBody()
	err := ffjson.Unmarshal(body, &data)
	if err != nil {
		return nil
	}
	return &data
}

func processEntityUpdate(ctx *fasthttp.RequestCtx, entity Entity) {
	data := loadJSON(ctx)
	if data != nil {
		ok := entity.Update(data, true)
		if ok {
			renderEmpty(ctx)
			return
		}
	}
	render400(ctx)
}

func processEntityCreate(ctx *fasthttp.RequestCtx, repo EntityRepo) {
	data := loadJSON(ctx)
	if data != nil {
		ok := repo.Create(data)
		if ok {
			renderEmpty(ctx)
			return
		}
	}
	render400(ctx)
}

func renderEntity(ctx *fasthttp.RequestCtx, entity Entity) {
	entity.to_json(ctx)
	ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}

func render400(ctx *fasthttp.RequestCtx) {
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
	ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}

func render404(ctx *fasthttp.RequestCtx) {
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}

func renderEmpty(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("{}"))
	ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}
