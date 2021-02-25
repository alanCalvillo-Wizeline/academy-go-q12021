package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("posts")
}

func SavePokemon(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
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
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		json.NewEncoder(w).Encode("Error while reading the File")
		log.Fatal(err)
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	//json.NewEncoder(w).Encode(params)
}
