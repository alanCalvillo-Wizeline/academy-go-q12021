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
var TotalPerWorker [3]int

const multiformSize = (10 << 20)

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Err := json.NewEncoder(w).Encode("posts")
	if Err != nil {
		log.Println(Err)
	}
}
func Read_csv(w http.ResponseWriter, r *http.Request) {
	errParse := r.ParseMultipartForm(multiformSize)
	if errParse != nil {
		log.Println(errParse)
	}
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		ErrEncode := json.NewEncoder(w).Encode("Error Retrieving the File")
		if ErrEncode != nil {
			log.Println(ErrEncode)
		}
		return
	}
	defer file.Close()
	var fs []Pokemon
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		CsvLine := strings.Split(scanner.Text(), ",")
		if len(CsvLine) == 2 {
			i2, err := strconv.ParseInt(CsvLine[0], 10, 64)
			if err != nil {
				log.Println("ID not being valid!")
			} else {
				a := Pokemon{Name: CsvLine[1], Id: int(i2)}
				fs = append(fs, a)
				models.SavePokemon(a.Id, a.Name)
			}

		} else {
			log.Println("CSV row incomplete!")
		}
	}

	if err := scanner.Err(); err != nil {
		ErrEncode := json.NewEncoder(w).Encode("Error while reading the File")
		if ErrEncode != nil {
			log.Println(ErrEncode)
		}
		log.Println(err)
	} else {
		j, _ := json.Marshal(fs)
		log.Println(string(j))
		ErrOutput := json.NewEncoder(w).Encode(fs)
		if ErrOutput != nil {
			log.Println(ErrOutput)
		}
		log.Println("Uploaded File: " + handler.Filename)
	}

}

func Api_feed(w http.ResponseWriter, r *http.Request) {
	// in this process we are going to use this url https://pokeapi.co/api/v2/pokemon?offset=0&limit=100
	ApiEndpoint := "https://pokeapi.co/api/v2/pokemon?offset=0&limit=100"
	response, err := http.Get(ApiEndpoint)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var ApiEndpointElement ApiEndpointBody
		json.Unmarshal(data, &ApiEndpointElement)
		var fs []Pokemon

		for i := 0; i < len(ApiEndpointElement.Results); i++ {
			url := ApiEndpointElement.Results[i].URL
			url = strings.Replace(url, "https://pokeapi.co/api/v2/pokemon/", "", 1)
			id := strings.Replace(url, "/", "", 1)
			id_poke, err := strconv.Atoi(id)
			if err != nil {
				fmt.Printf("The id extraction has been failed  %s\n", err)
			}
			a := Pokemon{Name: ApiEndpointElement.Results[i].Name, Id: id_poke}
			fs = append(fs, a)
			models.SavePokemon(a.Id, a.Name)
		}
		ErrEncode := json.NewEncoder(w).Encode(fs)
		if ErrEncode != nil {
			log.Println(ErrEncode)
		}
	}
}
func Worker_add(id int, jobs <-chan Pokemon, results chan<- Pokemon) {
	for j := range jobs {
		if total > 0 && TotalPerWorker[id-1] > 0 {
			fmt.Println("worker", id, "started  job", j)
			fmt.Println("worker", id, "finished job", j)
			results <- j
		}
		total--
		TotalPerWorker[id-1] = TotalPerWorker[id-1] - 1
	}
}
func Worker(w http.ResponseWriter, r *http.Request) {
	jobs := make(chan Pokemon, 100)
	results := make(chan Pokemon, 100)

	ItemsPerWorkers := r.URL.Query().Get("items_per_workers")
	IntValuePerWorker, _ := strconv.Atoi(ItemsPerWorkers)

	TypeValue := r.URL.Query().Get("type") //odd or even

	items := r.URL.Query().Get("items")

	ItemsValue, _ := strconv.Atoi(items)

	total = ItemsValue
	log.Println(ItemsPerWorkers)
	for w := 1; w <= 3; w++ {
		go Worker_add(w, jobs, results)
		TotalPerWorker[w-1] = IntValuePerWorker
	}
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		ErrEncode := json.NewEncoder(w).Encode("Error Retrieving the File")
		if ErrEncode != nil {
			log.Println(ErrEncode)
		}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		CsvLine := strings.Split(scanner.Text(), ",")
		if len(CsvLine) == 2 {
			i2, err := strconv.ParseInt(CsvLine[0], 10, 64)
			if err != nil {
				log.Println("ID not being valid!")
			} else {
				a := Pokemon{Name: CsvLine[1], Id: int(i2)}
				models.SavePokemon(a.Id, a.Name)
				if TypeValue == "even" && int(i2)%2 == 0 {
					jobs <- a
				} else {
					jobs <- a
				}

			}

		} else {
			log.Println("CSV row incomplete!")
		}
	}
	close(jobs)

	if err := scanner.Err(); err != nil {
		ErrEncode := json.NewEncoder(w).Encode("Error while reading the File")
		if ErrEncode != nil {
			log.Println(ErrEncode)
		}
		log.Println(err)
	} else {
		j, _ := json.Marshal(results)
		log.Println(string(j))
		ErrEncode := json.NewEncoder(w).Encode(results)
		if ErrEncode != nil {
			log.Println(ErrEncode)
		}
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
		ErrEncode := json.NewEncoder(w).Encode("not supported")
		if ErrEncode != nil {
			log.Println(ErrEncode)
		}
	}

}
