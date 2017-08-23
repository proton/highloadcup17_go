package main

import (
	"fmt"
	//"github.com/pquerna/ffjson/ffjson"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func startWebServer() {
	fmt.Println("Webserver: starting")
	// h := requestHandler
	h := timeoutHandler
	// h = fasthttp.CompressHandler(h)

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func timeoutHandler(ctx *fasthttp.RequestCtx) {
	doneCh := make(chan struct{})
	go func() {
		requestHandler(ctx)
		close(doneCh)
	}()

	select {
	case <-doneCh:
		// fmt.Println("The task has been finished in less than a second")
	case <-time.After(time.Second):
		fmt.Println("Timeout")
		fmt.Printf("\n\nWEB SERVER ERROR: %s %s - %s\n%s\n", string(ctx.Method()), string(ctx.Path()), string(ctx.PostBody()), debug.Stack())
		ctx.TimeoutError("Timeout!")
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	// ctx.SetContentType("text/plain; charset=utf8")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n\nWEB SERVER ERROR: %s %s - %s\n%s\n%s\n", string(ctx.Method()), string(ctx.Path()), string(ctx.PostBody()), r, debug.Stack())
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
			entity, ok := find_entity(repo, &path[2])
			if ok {
				renderEntity(ctx, entity)
				return
			}
		} else if path[2] == "new" {
			processEntityCreate(ctx, repo)
			return
		} else {
			entity, ok := find_entity(repo, &path[2])
			if ok {
				processEntityUpdate(ctx, entity)
				return
			}
		}
	} else if path_len == 4 && http_method_is_get {
		if entity_kind == "users" && path[3] == "visits" {
			user, ok := find_user(&path[2])
			if ok {
				processUserVisits(ctx, user)
				return
			}
		} else if entity_kind == "locations" && path[3] == "avg" {
			location, ok := find_location(&path[2])
			if ok {
				processLocationAvgs(ctx, location)
				return
			}
		}
	}
	render404(ctx)
}

func extractUintParam(ctx *fasthttp.RequestCtx, key string) (*int, bool) {
	param := ctx.QueryArgs().Peek(key)
	if param == nil {
		return nil, true
	}
	param_int, err := strconv.Atoi(string(param))
	if err != nil {
		return nil, false
	}
	param_uint := int(param_int)
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
	if strings.Contains(string(body), ": null") {
		return nil
	}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil
	}
	return &data
}

func processEntityUpdate(ctx *fasthttp.RequestCtx, entity Entity) {
	data := loadJSON(ctx)
	if data == nil {
		render400(ctx)
		return
	}
	entity.Update(data, true)
	renderEmpty(ctx)
	// ctx.SetConnectionClose() // is it really helps?
	// go entity.Update(data, true)
}

func processEntityCreate(ctx *fasthttp.RequestCtx, repo EntityRepo) {
	data := loadJSON(ctx)
	if data == nil {
		render400(ctx)
		return
	}
	repo.Create(data)
	renderEmpty(ctx)
	// ctx.SetConnectionClose() // is it really helps?
	// go repo.Create(data)
}

func renderEntity(ctx *fasthttp.RequestCtx, entity Entity) {
	entity.toJson(ctx)
	// ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}

func render400(ctx *fasthttp.RequestCtx) {
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
	// ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}

func render404(ctx *fasthttp.RequestCtx) {
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	// ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}

func renderEmpty(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("{}"))
	// ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}
