package main

import (
	"bytes"
	"fmt"
	//"github.com/pquerna/ffjson/ffjson"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	// "runtime/debug"
	"strconv"
	// "time"
)

func startWebServer() {
	fmt.Println("Webserver: starting")
	h := requestHandler
	// h := timeoutHandler
	// h = fasthttp.CompressHandler(h)

	if err := fasthttp.ListenAndServe(*ADDR, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

// func timeoutHandler(ctx *fasthttp.RequestCtx) {
// 	doneCh := make(chan struct{})
// 	go func() {
// 		requestHandler(ctx)
// 		close(doneCh)
// 	}()

// 	select {
// 	case <-doneCh:
// 		// fmt.Println("The task has been finished in less than a second")
// 	case <-time.After(time.Second):
// 		fmt.Println("Timeout")
// 		fmt.Printf("\n\nWEB SERVER ERROR: %s %s - %s\n%s\n", string(ctx.Method()), string(ctx.Path()), string(ctx.PostBody()), debug.Stack())
// 		ctx.TimeoutError("Timeout!")
// 	}
// }

var (
	METHOD_GET    = []byte("GET")
	PATH_SPLITTER = []byte("/")
	B_USERS       = []byte("users")
	B_LOCATIONS   = []byte("locations")
	B_VISITS      = []byte("visits")
	B_NEW         = []byte("new")
	B_AVG         = []byte("avg")
	NULL_CHECK    = []byte(": null")
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	// ctx.SetContentType("text/plain; charset=utf8")
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("\n\nWEB SERVER ERROR: %s %s - %s\n%s\n%s\n", string(ctx.Method()), string(ctx.Path()), string(ctx.PostBody()), r, debug.Stack())
	// 		render400(ctx)
	// 	}
	// }()

	http_method_is_get := bytes.Equal(ctx.Method(), METHOD_GET)
	path := bytes.Split(ctx.Path(), PATH_SPLITTER)
	path_len := len(path)
	repo := entity_repo(len(path[1]))
	if path_len == 3 && !http_method_is_get && bytes.Equal(path[2], B_NEW) {
		processEntityCreate(ctx, repo)
		return
	}

	entity_id_str := string(path[2])

	if path_len == 3 {
		entity, ok := find_entity(repo, &entity_id_str)
		if ok {
			if http_method_is_get {
				renderEntity(ctx, entity)
			} else {
				processEntityUpdate(ctx, entity)
			}
			return
		}
	} else if path_len == 4 && http_method_is_get {
		if bytes.Equal(path[3], B_VISITS) {
			user, ok := find_user(&entity_id_str)
			if ok {
				processUserVisits(ctx, user)
				return
			}
		} else if bytes.Equal(path[3], B_AVG) {
			location, ok := find_location(&entity_id_str)
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
	if bytes.Contains(body, NULL_CHECK) {
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
	// entity.Update(data, true)
	renderEmpty(ctx)
	// ctx.SetConnectionClose() // is it really helps?
	go entity.Update(data, true)
}

func processEntityCreate(ctx *fasthttp.RequestCtx, repo EntityRepo) {
	data := loadJSON(ctx)
	if data == nil {
		render400(ctx)
		return
	}
	// repo.Create(data)
	renderEmpty(ctx)
	// ctx.SetConnectionClose() // is it really helps?
	go repo.Create(data)
}

func renderEntity(ctx *fasthttp.RequestCtx, entity Entity) {
	entity.writeJSON(ctx)
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
	ctx.WriteString("{}")
	// ctx.SetConnectionClose() // https://github.com/sat2707/hlcupdocs/issues/37
}
