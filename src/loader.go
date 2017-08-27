package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func unpackDataZip() {
	file_path := *DATAZIP_PATH
	if _, err := os.Stat(file_path); !os.IsNotExist(err) {
		fmt.Println("DataLoading: extracting zip file")
		cmd := exec.Command("unzip", file_path, "-d", *DATA_DIR)
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func loadInitialTime() {
	file_path := *OPTIONS_PATH
	data, _ := ioutil.ReadFile(file_path)
	ts_str := strings.Split(string(data), "\n")[0]
	ts, _ := strconv.Atoi(ts_str)
	InitialTime = time.Unix(int64(ts), 0)
	fmt.Println("DataLoading: set timestamp to", InitialTime)
}

func loadJsons() {
	var data []byte

	files, _ := ioutil.ReadDir(*DATA_DIR)

	entity_kinds := []string{"users", "locations", "visits"}

	for _, entity_kind := range entity_kinds {
		repo := entity_repo(len(entity_kind))

		for _, f := range files {
			if !strings.Contains(f.Name(), entity_kind) {
				continue
			}
			fmt.Println("DataLoading: loading", f.Name())

			file_path := *DATA_DIR + f.Name()

			data, _ = ioutil.ReadFile(file_path)

			jsonparser.ArrayEach(data, func(object_data []byte, dataType jsonparser.ValueType, offset int, err error) {
				repo.CreateFromJSON(object_data)
			}, entity_kind)
		}
	}
}

func loadInitialData() {
	fmt.Println("DataLoading: starting")
	unpackDataZip()
	loadInitialTime()
	loadJsons()
	fmt.Println("DataLoading: finished")
}
