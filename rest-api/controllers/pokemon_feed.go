package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rest-api/models"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

/*
	library to be used for pokemons : https://pokeapi.co/api/v2/pokemon/{poke-id}
*/
type Pokemon struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ApiEndpointBody struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var total int
var total_per_worker [3]int

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("posts")
}
func Read_csv(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Fatal("Error Retrieving the File")
		json.NewEncoder(w).Encode("Error Retrieving the File")
		return
	}
	defer file.Close()
	var fs []Pokemon
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		csv_line := strings.Split(scanner.Text(), ",")
		if len(csv_line) == 2 {
			i2, err := strconv.ParseInt(csv_line[0], 10, 64)
			if err != nil {
				log.Fatal("ID not being valid!")
			} else {
				a := Pokemon{Name: csv_line[1], Id: int(i2)}
				fs = append(fs, a)
				models.SavePokemon(a.Id, a.Name)
			}

		} else {
			log.Fatal("CSV row incomplete!")
		}
	}

	if err := scanner.Err(); err != nil {
		json.NewEncoder(w).Encode("Error while reading the File")
		log.Fatal(err)
	} else {
		j, _ := json.Marshal(fs)
		log.Println(string(j))
		json.NewEncoder(w).Encode(fs)
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	}

}

func Api_feed(w http.ResponseWriter, r *http.Request) {
	// in this process we are going to use this url https://pokeapi.co/api/v2/pokemon?offset=0&limit=100
	api_endpoint := "https://pokeapi.co/api/v2/pokemon?offset=0&limit=100"
	response, err := http.Get(api_endpoint)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var api_endpoint_body ApiEndpointBody
		json.Unmarshal(data, &api_endpoint_body)
		var fs []Pokemon

		for i := 0; i < len(api_endpoint_body.Results); i++ {
			url := api_endpoint_body.Results[i].URL
			url = strings.Replace(url, "https://pokeapi.co/api/v2/pokemon/", "", 1)
			id := strings.Replace(url, "/", "", 1)
			id_poke, err := strconv.Atoi(id)
			if err != nil {
				fmt.Printf("The id extraction has been failed  %s\n", err)
			}
			a := Pokemon{Name: api_endpoint_body.Results[i].Name, Id: id_poke}
			fs = append(fs, a)
			models.SavePokemon(a.Id, a.Name)
		}
		json.NewEncoder(w).Encode(fs)
	}
}
func worker_add(id int, jobs <-chan Pokemon, results chan<- Pokemon) {
	for j := range jobs {
		if total > 0 && total_per_worker[id-1] > 0 {
			fmt.Println("worker", id, "started  job", j)
			fmt.Println("worker", id, "finished job", j)
			results <- j
		}
		total--
		total_per_worker[id-1] = total_per_worker[id-1] - 1
	}
}
func Worker(w http.ResponseWriter, r *http.Request) {
	jobs := make(chan Pokemon, 100)
	results := make(chan Pokemon, 100)

	items_per_workers := r.URL.Query().Get("items_per_workers")
	intValue_per_worker, _ := strconv.Atoi(items_per_workers)

	type_value := r.URL.Query().Get("type") //odd or even

	items := r.URL.Query().Get("items")

	items_value, _ := strconv.Atoi(items)

	total = items_value
	fmt.Printf(items_per_workers)
	for w := 1; w <= 3; w++ {
		go worker_add(w, jobs, results)
		total_per_worker[w-1] = intValue_per_worker
	}
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Fatal("Error Retrieving the File")
		json.NewEncoder(w).Encode("Error Retrieving the File")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		csv_line := strings.Split(scanner.Text(), ",")
		if len(csv_line) == 2 {
			i2, err := strconv.ParseInt(csv_line[0], 10, 64)
			if err != nil {
				log.Fatal("ID not being valid!")
			} else {
				a := Pokemon{Name: csv_line[1], Id: int(i2)}
				models.SavePokemon(a.Id, a.Name)
				if type_value == "even" && int(i2)%2 == 0 {
					jobs <- a
				} else {
					jobs <- a
				}

			}

		} else {
			log.Fatal("CSV row incomplete!")
		}
	}
	close(jobs)

	if err := scanner.Err(); err != nil {
		json.NewEncoder(w).Encode("Error while reading the File")
		log.Fatal(err)
	} else {
		j, _ := json.Marshal(results)
		log.Println(string(j))
		json.NewEncoder(w).Encode(results)
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)

	}

}

func SavePokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	switch params["method"] {
	case "csv":
		Read_csv(w, r)
	case "api":
		Api_feed(w, r)
	case "worker":
		Worker(w, r)
	default:
		json.NewEncoder(w).Encode("not supported")
	}

}
