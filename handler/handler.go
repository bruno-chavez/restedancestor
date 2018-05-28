// Package handler is used to separate the handlers from other functions of the API.
package handler

import (
	"encoding/json"
	"github.com/bruno-chavez/restedancestor/database"
	"net/http"
)

// slice is a global variable to avoid multiple calls to Parser since always returns the same slice.
var slice = database.Parser()

// RandomHandler takes care of the 'random' route.
func RandomHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		marshaledData, _ := json.MarshalIndent(slice.Random(), "", "")

		w.Header().Set("Content-Type", "application/json")
		w.Write(marshaledData)

	case "OPTIONS":
		w.Header().Set("Allow", "GET,OPTIONS")
	}

}

// AllHandler takes care of the 'all' route.
func AllHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		marshaledData, _ := json.MarshalIndent(slice, "", "")

		w.Header().Set("Content-Type", "application/json")
		w.Write(marshaledData)
	case "OPTIONS":
		w.Header().Set("Allow", "GET,OPTIONS")
	}
}
