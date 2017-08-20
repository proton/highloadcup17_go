package main

// https://github.com/sat2707/hlcupdocs

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	// "github.com/hashicorp/go-memdb"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"strings"
	// "time"
	"runtime/debug" //TODO:
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

	for _, f := range r.File {
		if !strings.Contains(f.Name, ".json") {
			continue
		}
		fmt.Println("DataLoading: loading", f.Name)
		arr := strings.Split(f.Name, "_")
		entity_kind := arr[0]

		rc, _ := f.Open()
		b, _ := ioutil.ReadAll(rc)
		data := make(JsonDataArray)
		json.Unmarshal(b, &data)
		json_objects := data[entity_kind]

		repo := entity_repo(entity_kind)
		for _, json_object := range json_objects {
			repo.Create(&json_object)
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
	startWebServer()
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain; charset=utf8")
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(ctx, "\n\nWEB SERVER ERROR: %s\n%s\n", r, debug.Stack())
		}
	}()

	path := strings.Split(string(ctx.Path()), "/")
	http_method := string(ctx.Method())
	path_len := len(path)
	entity_kind := path[1]
	repo := entity_repo(entity_kind)

	if path_len == 3 {
		var entity, ok = find_entity(repo, &path[2])
		if ok == true {
			if http_method == "GET" {
				renderEntity(ctx, entity)
			} else {
				processEntityUpdate(ctx, entity)
			}
			return
		}
	} else if path_len == 4 {
		if http_method == "GET" {
			if entity_kind == "users" && path[3] == "visits" {
				var user, ok = find_user(&path[2])
				if ok == true {
					processUserVisits(ctx, user)
					return
				}
			}
			if entity_kind == "locations" && path[3] == "avg" {
				// TODO:
			}
		} else {
			if path[3] == "new" {
				processEntityCreate(ctx, repo)
				return
			}
		}
	}
	render404(ctx)
}

func find_entity(repo EntityRepo, entity_id_str *string) (Entity, bool) {
	entity_id_int, error := strconv.Atoi(*entity_id_str)
	if error == nil {
		entity_id := uint32(entity_id_int)
		return repo.FindEntity(entity_id)
	}
	return nil, false
}

func find_user(entity_id_str *string) (*User, bool) {
	entity_id_int, error := strconv.Atoi(*entity_id_str)
	if error == nil {
		entity_id := uint32(entity_id_int)
		return Users.Find(entity_id)
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

func processUserVisits(ctx *fasthttp.RequestCtx, user *User) {
	fromDate, ok := extractUintParam(ctx, "fromDate")
	if ok == false {
		render400(ctx)
		return
	}
	toDate, ok := extractUintParam(ctx, "toDate")
	if ok == false {
		render400(ctx)
		return
	}
	country, ok := extractStringParam(ctx, "country")
	if ok == false {
		render400(ctx)
		return
	}
	toDistance, ok := extractUintParam(ctx, "toDistance")
	if ok == false {
		render400(ctx)
		return
	}

	user.WriteVisitsJson(ctx, fromDate, toDate, country, toDistance)

	// fmt.Fprintf(ctx, "myVariable = %#v \n", arg_fromDate)
	// if arg_fromDate == nil {
	// 	fmt.Fprintf(ctx, "arr\n")
	// }

	// query_string := string(ctx.QueryArgs())
	// query_args := strings.Split(query_string, "&")
	// for _, query_arg := range r.File {
	// 	query_key, query_string := strings.Split(query_string, "=")
	// 	if len(query_string) == 0 {
	// 		render400(ctx)
	// 		return
	// 	}
	// }

	// fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	// fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	// fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	// fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	// fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	// fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	// fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	// fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	// fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())
	// fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)
}

func processEntityUpdate(ctx *fasthttp.RequestCtx, entity Entity) {
	request_body := ctx.PostBody()
	var data, err = readRequstJson(request_body)
	if err != nil {
		render400(ctx)
	} else {
		entity.Update(&data, true)
		renderEmpty(ctx)
	}
}

func processEntityCreate(ctx *fasthttp.RequestCtx, repo EntityRepo) {
	request_body := ctx.PostBody()
	err := repo.CreateFromJson(request_body)
	if err != nil {
		render400(ctx)
		return
	}
	renderEmpty(ctx)
}

func renderEntity(ctx *fasthttp.RequestCtx, entity Entity) {
	entity.to_json(ctx)
}

func render400(ctx *fasthttp.RequestCtx) {
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
}

func render404(ctx *fasthttp.RequestCtx) {
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
}

func renderEmpty(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("{}"))
}
