package controller

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
  )

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("posts")
  }

func SavePokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	
	json.NewEncoder(w).Encode(params)
  }