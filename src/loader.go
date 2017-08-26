package main

import (
	"archive/zip"
	"fmt"
	"github.com/buger/jsonparser"
	// "github.com/pquerna/ffjson/ffjson"
	// "encoding/json"
	// "github.com/pkg/profile"
	"io/ioutil"
	"log"
	"strings"
)

func loadInitialData() {
	fmt.Println("DataLoading: starting")
	r, err := zip.OpenReader(*DATA_PATH)
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

			repo := entity_repo(len(entity_kind))

			rc, _ := f.Open()
			data, _ := ioutil.ReadAll(rc)

			jsonparser.ArrayEach(data, func(object_data []byte, dataType jsonparser.ValueType, offset int, err error) {
				repo.CreateFromJSON(object_data)
			}, entity_kind)
		}
	}
	fmt.Println("DataLoading: finished")
}
