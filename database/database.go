// Package database takes care of properly handle the database to be used in other parts of the API.
package database

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

// init is used to seed the rand.Intn function.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// QuoteType is used to parse the whole json database in a slice of the QuoteType type.
type QuoteType struct {
	Quote string    `json:"quote"`
	Uuid  uuid.UUID `json:"uuid"`
	Score int       `json:"score"`
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
	path := os.Getenv("GOPATH") + "/src/github.com/bruno-chavez/restedancestor/database/database.json"
	goingBack := ""

	path = goingBack + path

	rawJSON, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer rawJSON.Close()

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


		log.Fatal(err)
	}
}

func OffsetQuoteFromUUID(uuid string) (*int, error) {
	path := os.Getenv("GOPATH") + "/src/github.com/bruno-chavez/restedancestor/database/database.json"
	goingBack := ""

	path = goingBack + path

	rawJSON, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer rawJSON.Close()

	readJSON, err2 := ioutil.ReadAll(rawJSON)
	if err2 != nil {
		log.Fatal(err2)
	}

	parsedJSON := make(QuoteSlice, 0)
	err3 := json.Unmarshal(readJSON, &parsedJSON)
	if err3 != nil {
		log.Fatal(err3)
	}

	for k, quote := range parsedJSON {
		if quote.Uuid.String() == uuid {
			return &k, nil
		}
	}

	return nil, errors.New("Unknown")
}

// LikeQuote increments the score of the quote
func LikeQuote(uuid string) {
	offset, _ := OffsetQuoteFromUUID(uuid)
	slice := Parser()
	slice[*offset].Score += 1

	writeJSON, err4 := json.MarshalIndent(slice, "", "  ")
	if err4 != nil {
		log.Fatal(err4)
	}

	path := os.Getenv("GOPATH") + "/src/github.com/bruno-chavez/restedancestor/database/database.json"
	goingBack := ""

	path = goingBack + path
	ioutil.WriteFile(path, writeJSON, 0)
}

// DislikeQuote decrements the score of the quote
func DislikeQuote(uuid string) {
	offset, _ := OffsetQuoteFromUUID(uuid)
	slice := Parser()
	slice[*offset].Score -= 1

	writeJSON, err4 := json.MarshalIndent(slice, "", "  ")
	if err4 != nil {
		log.Fatal(err4)
	}

	path := os.Getenv("GOPATH") + "/src/github.com/bruno-chavez/restedancestor/database/database.json"
	goingBack := ""

	path = goingBack + path
	ioutil.WriteFile(path, writeJSON, 0)
}
