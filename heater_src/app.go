package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	flag.Parse()
	log.Println("Heater: started")

	var wg sync.WaitGroup
	entity_kinds := []string{"users", "locations", "visits"}
	wg.Add(len(entity_kinds))
	for _, entity_kind := range entity_kinds {
		go func(entity_kind string) {
			defer wg.Done()
			for id := 1; id <= 10000; id++ {
				// log.Println("Heater: it", it)
				url := makeUrl(entity_kind, id)
				makeRequest(url)
				time.Sleep(10 * time.Millisecond)
			}
		}(entity_kind)
	}

	wg.Wait()

	log.Println("Heater: finished")
}

func makeUrl(entity_kind string, id int) string {
	return "http://localhost" + *PORT + "/" + entity_kind + "/" + strconv.Itoa(id)
}

func makeRequest(url string) {
	// log.Println("Heater: requesting", url)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	} else {
		defer res.Body.Close()
		_, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		} else {
			// log.Println("Heater: requesting", url, ":", string(body))
		}
	}
}

var (
	PORT = flag.String("addr", ":80", "TCP address to send requests")
)
