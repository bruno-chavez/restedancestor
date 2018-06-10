// Package database takes care of properly handle the database to be used in other parts of the API.
package database

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

// init is used to seed the rand.Intn function.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// QuoteType is used to parse the whole json database in a slice of the QuoteType type.
type QuoteType struct {
	Quote string `json:"quote"`
}

// QuoteSlice exists to provide abstraction to the QuoteType type,
// since its always going to be used as a slice.
type QuoteSlice []QuoteType

// Random returns a random quote from a QuoteSlice type.
func (q QuoteSlice) Random() QuoteType {
	return q[rand.Intn(len(q))]
}

// Parser fetches from database.json and puts it on a slice.
func Parser() QuoteSlice {
	path := "database/database.json"

	readJSON, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	parsedJSON := make(QuoteSlice, 0)
	err2 := json.Unmarshal(readJSON, &parsedJSON)
	if err2 != nil {
		panic(err2)
	}

	return parsedJSON
}
