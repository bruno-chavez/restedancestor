// Package main is the start of the API, creates all the routes and sets a handler for each one.
package main

import (
	"fmt"
	"github.com/bruno-chavez/restedancestor/handler"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")

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
	fmt.Println("The Ancestor is waiting and listening on port ", port)
	http.ListenAndServe(":" + port, router)
}
