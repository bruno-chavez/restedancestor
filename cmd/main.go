package main

import (
	"fmt"
	"github.com/bruno-chavez/restedancestor/server"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	router := mux.NewRouter()
	router.NewRoute().
		Path("/rquote").
		HandlerFunc(server.RquoteHandler).
		Methods("GET")
	router.NewRoute().
		Path("/allquote").
		HandlerFunc(server.AllquoteHandler).
		Methods("GET")

	fmt.Println("Welcome to restedancestor, the API is running in a maddening fashion!")
	fmt.Println("The Ancestor is waiting and listening on port 8000 of localhost")
	http.ListenAndServe(":8000", router)
}
