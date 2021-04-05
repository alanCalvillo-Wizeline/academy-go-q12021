package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	controller "rest-api/controllers"
	"testing"
)

type PokeElement []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Test_worker_endpoint(t *testing.T) {
	// make a http request by using httptest
	req := httptest.NewRequest(http.MethodGet, "/pokemon/worker?items_per_workers=3&items=3&type=even", nil)
	// make response recorder for get response that send out
	rw := httptest.NewRecorder()
	// expected result
	expectedResult := false

	// call an api
	controller.Api_feed(rw, req)

	// response from api
	response := rw.Result()
	body, _ := ioutil.ReadAll(response.Body)

	var poke_elements PokeElement
	json.Unmarshal(body, &poke_elements)
	for i := 0; i < len(poke_elements); i++ {
		if poke_elements[i].ID == 1 && poke_elements[i].Name == "bulbasaur" {
			expectedResult = true
		}
	}
	// because body is []byte, we have to change to string to compare with expectedResult
	if expectedResult != true {
		t.Errorf("bad response")
	}
}
