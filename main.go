// Package main is the start of the API, creates all the routes and sets a handlers for each one.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bruno-chavez/restedancestor/handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// initiates router
	router := httprouter.New()

	// lists routes
	router.GET("/random", handlers.Random)
	router.GET("/all", handlers.All)
	router.GET("/senile", handlers.Senile)
	router.GET("/search/:word", handlers.Search)

	//uuid routes
	router.GET("/uuid/:uuid/find", handlers.Find)
	router.POST("/uuid/:uuid/like", handlers.Like)
	router.POST("/uuid/:uuid/dislike", handlers.Dislike)
	router.GET("/uuid/top", handlers.Top)

	fmt.Println("Welcome to restedancestor, the API is running in a maddening fashion!")
	fmt.Println("The Ancestor is waiting and listening on localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
