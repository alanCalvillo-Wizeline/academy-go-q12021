package main
import (
  "github.com/gorilla/mux"
  "net/http"
  "./controllers"
)

func main() {
  router := mux.NewRouter()
  
  router.HandleFunc("/pokemon", controller.GetPokemon).Methods("GET")
  /*
    method attribute will be used to check if we have a CSV save or External API save in order
    to map with the proper function in the controller
  */
  router.HandleFunc("/pokemon/{method}", controller.SavePokemon).Methods("POST")
  
http.ListenAndServe(":8000", router)
}