package quotes

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/bruno-chavez/restedancestor/database"
	"github.com/satori/go.uuid"
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
func Parser(data database.Database) QuoteSlice {
	parsedJSON := make(QuoteSlice, 0)
	err3 := json.Unmarshal(data.Read(), &parsedJSON)
	if err3 != nil {
		log.Fatal(err3)
	}

	return parsedJSON
}

// LikeQuote increments the score of the quote
func (q QuoteSlice) LikeQuote(db database.Database, uuid string) {
	offset, _ := q.OffsetQuoteFromUUID(db, uuid)
	q[*offset].Score++

	if err := unparser(db, q); err != nil {
		log.Fatal(err)
	}
}

// DislikeQuote decrements the score of the quote
func (q QuoteSlice) DislikeQuote(db database.Database, uuid string) {
	offset, _ := q.OffsetQuoteFromUUID(db, uuid)
	q[*offset].Score--

	if err := unparser(db, q); err != nil {
		log.Fatal(err)
	}
}

// OffsetQuoteFromUUID find the uuid in the slice and returns its offset
func (q QuoteSlice) OffsetQuoteFromUUID(db database.Database, uuid string) (*int, error) {

	for k, quote := range q {

		if quote.Uuid.String() == uuid {
			return &k, nil
		}
	}

	return nil, errors.New("unknown")
}

// unparser writes a slice into database.
func unparser(db database.Database, quotes QuoteSlice) error {
	writeJSON, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return db.Write(writeJSON)
}
