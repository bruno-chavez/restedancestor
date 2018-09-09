// Package database takes care of properly handle the database to be used in other parts of the API.
package database

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

// init is used to seed the rand.Intn function.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// NewDb initializes Database struct
func NewDb(p string) database {
	database := database{
		path: p,
	}

	return database
}

// QuoteType is used to parse the whole json database in a slice of the QuoteType type.
type QuoteType struct {
	Quote string `json:"quote"`
}

// QuoteSlice exists to provide abstraction to the QuoteType type,
// since its always going to be used as a slice.
type QuoteSlice []QuoteType

// database is a data structure embedding all behaviors involving storage
type database struct {
	path  string
	slice QuoteSlice
}

// Random returns a random quote from a QuoteSlice type.
func (q QuoteSlice) Random() QuoteType {
	return q[rand.Intn(len(q))]
}

// Parser fetches from database.json and puts it on a slice.
func (d database) Parser() QuoteSlice {
	if d.slice == nil {
		rawJSON, err := os.Open(d.path)
		if err != nil {
			log.Fatal(err)
		}

		readJSON, err2 := ioutil.ReadAll(rawJSON)
		if err2 != nil {
			log.Fatal(err2)
		}

		parsedJSON := make(QuoteSlice, 0)
		err3 := json.Unmarshal(readJSON, &parsedJSON)
		if err3 != nil {
			log.Fatal(err3)
		}
		d.slice = parsedJSON
	}
	return d.slice
}

// Path returns the storage path
func (d database) Path() string {
	return d.path
}
