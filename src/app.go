package main

// https://github.com/sat2707/hlcupdocs

import (
	"archive/zip"
	"encoding/json"
	"flag"
	// "fmt"
	// "github.com/hashicorp/go-memdb"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"strings"
	// "time"
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
	r, err := zip.OpenReader("/tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if !strings.Contains(f.Name, ".json") {
			continue
		}
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

	path := strings.Split(string(ctx.Path()), "/")
	http_method := string(ctx.Method())
	path_len := len(path)
	entity_kind := path[1]
	var entity_id_int, _ = strconv.Atoi(path[2])
	entity_id := uint32(entity_id_int)

	if path_len == 3 {
		repo := entity_repo(entity_kind)
		var entity, ok = repo.Find(entity_id)
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
				// TODO:
			}
			if entity_kind == "locations" && path[3] == "avg" {
				// TODO:
			}
		} else {
			if path[3] == "new" {
				repo := entity_repo(entity_kind)
				processEntityCreate(ctx, repo)
				return
			}
		}
	}
	render404(ctx)

	// fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
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
		entity.Update(&data)
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
