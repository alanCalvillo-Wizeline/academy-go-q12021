package main

import "fmt"

func hello(user string) string {
	if ( len(user) == 0){
		return "Hello"
	} else {
		return "Hello "+user
	}
}
func main() {
	greeting := hello("Alan")
	fmt.Println(greeting)
}