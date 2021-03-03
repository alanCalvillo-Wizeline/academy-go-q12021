package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("posts")
}
func read_csv(w http.ResponseWriter, r *http.Request) {
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

func SavePokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	switch params["method"] {
	case "csv":
		read_csv(w, r)
	case "api":
		json.NewEncoder(w).Encode("it will work at some point")
	default:
		json.NewEncoder(w).Encode("not supported")
	}

}
