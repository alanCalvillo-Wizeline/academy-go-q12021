package main
import "testing"


func TestHello(t *testing.T) {
	emptyResult := hello("Edgar")

	if emptyResult == "Hello" {
		t.Error("hello(\"\") failed")
	} else {
		t.Log("success")
	}
	
	result := hello("Alan")

	if result != "Hello Alan" {
		t.Error("failed")
	} else {
		t.Log("success")
	}
	
}