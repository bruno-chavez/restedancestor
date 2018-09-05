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

func (q QuoteSlice) Len() int {
	return len(q)
}

func (q QuoteSlice) Swap(i int, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q QuoteSlice) Less(i int, j int) bool {
	return q[i].Score > q[j].Score
}

// Parser fetches from database.json and puts it on a slice.
func Parser() QuoteSlice {
	rawJSON, err := os.Open(path())
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

// LikeQuote increments the score of the quote
func LikeQuote(uuid string) {
	offset, _ := OffsetQuoteFromUUID(uuid)
	slice := Parser()
	slice[*offset].Score += 1

	if err := unparser(slice); err != nil {
		log.Fatal(err)
	}
}

// DislikeQuote decrements the score of the quote
func DislikeQuote(uuid string) {
	offset, _ := OffsetQuoteFromUUID(uuid)
	slice := Parser()
	slice[*offset].Score -= 1

	if err := unparser(slice); err != nil {
		log.Fatal(err)
	}
}

// OffsetQuoteFromUUID find the uuid in the slice and returns its offset
func OffsetQuoteFromUUID(uuid string) (*int, error) {
	parsedJSON := Parser()

	for k, quote := range parsedJSON {
		if quote.Uuid.String() == uuid {
			return &k, nil
		}
	}

	return nil, errors.New("Unknown")
}

// unparser writes a slice into database.json.
func unparser(quotes QuoteSlice) error {
	writeJSON, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return ioutil.WriteFile(path(), writeJSON, 0)
}

func path() string {
	return os.Getenv("GOPATH") + "/src/github.com/bruno-chavez/restedancestor/database/database.json"
}
