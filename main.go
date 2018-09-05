// Package main is the start of the API, creates all the routes and sets a handler for each one.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bruno-chavez/restedancestor/handler"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.NewRoute().
		Path("/random").
		HandlerFunc(handler.RandomHandler).
		Methods("GET", "OPTIONS")
	router.NewRoute().
		Path("/all").
		HandlerFunc(handler.AllHandler).
		Methods("GET", "OPTIONS")
	router.NewRoute().
		Path("/search/{word}").
		HandlerFunc(handler.SearchHandler).
		Methods("GET", "OPTIONS")

	fmt.Println("Welcome to restedancestor, the API is running in a maddening fashion!")
	fmt.Println("The Ancestor is waiting and listening on localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
