// Package handler is used to separate the handlers from other functions of the API.
package handler

import (
	"encoding/json"
	"github.com/bruno-chavez/restedancestor/database"
	"github.com/bruno-chavez/restedancestor/lib"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
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

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		query := mux.Vars(r)
		wordToSearch := query["word"]

		matched := false
		quotesMatched := make(database.QuoteSlice, 0)

		// quote is a QuoteType type.
		// xSlice is a slice with one word and white-spaces that appear after a filter is applied.
		// xFilter is a string inside a wordXSlice.
		// This for is a candidate for been transformed into an recursive function inside lib.
		for _, quote := range slice {
			firstSlice := strings.Split(quote.Quote, " ")

			for _, firstFilter := range firstSlice {
				secondSlice := strings.Split(firstFilter, ",")

				for _, secondFilter := range secondSlice {
					thirdSlice := strings.Split(secondFilter, "!")

					for _, thirdFilter := range thirdSlice {
						filteredWord := strings.Split(thirdFilter, ".")

						// After all the filters are applied the filtered word is compared with the word that is been searched.
						// If a match is found, the quote that the filtered word belongs to, is printed.
						if filteredWord[0] == wordToSearch {
							quotesMatched = append(quotesMatched, quote)
							matched = true
						}
					}
				}
			}
		}
		if matched {
			w.Header().Set("Content-Type", "application/json")
			filteredSLice, _ := json.MarshalIndent(quotesMatched, "", "")
			w.Write(filteredSLice)

		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			notfoundJSON := lib.NotFound(wordToSearch)
			w.Write(notfoundJSON)
		}

	case "OPTIONS":
		w.Header().Set("Allow", "GET,OPTIONS")
	}
}
