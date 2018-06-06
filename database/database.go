// Package database takes care of properly handle the database to be used in other parts of the API.
package database

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"log"
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

	// Don't worry too much about how Parser works, trust me, it does!
	currentDir, _ := os.Getwd()
	totalDoubleDots := len(strings.Split(currentDir, "/"))
	path := os.Getenv("GOPATH") + "/src/github.com/bruno-chavez/restedancestor/database/database.json"
	goingBack := ""
	for i := 1; i <= totalDoubleDots; i++ {
		if i == totalDoubleDots {
			goingBack = goingBack + ".."
		} else {
			goingBack = "../" + goingBack
		}
	}
	path = goingBack + path

	rawJSON, err := os.Open(path)
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

	return parsedJSON
}
