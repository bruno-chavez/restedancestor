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

// QuoteType describes a quote.
type QuoteType struct {
	ID    int
	Quote string    `json:"quote"`
	Uuid  uuid.UUID `json:"uuid"`
	Score int       `json:"score"`
}

// Quotes is used to parse the whole json database.
type Quotes struct {
	Data []QuoteType `json:"data"`
}

// Random returns a random quote from a Quotes type.
func (q Quotes) Random() QuoteType {
	qd := q.Data
	return qd[rand.Intn(len(qd))]
}

// Parser fetches from database.json and puts it on a struct.
func Parser(data database.Storage) Quotes {
	q := Quotes{}
	err := json.Unmarshal(data.Read(), &q)
	if err != nil {
		log.Fatal(err)
	}

	return q
}

// OffsetQuoteFromUUID find the uuid in the slice and returns its offset
func (q Quotes) OffsetQuoteFromUUID(uuid string) (*int, error) {
	for k, quote := range q.Data {
		if quote.Uuid.String() == uuid {
			return &k, nil
		}
	}

	return nil, errors.New("unknown")
}

// List returns all quotes containing a word
// func (q Quotes) List(w string) []QuoteType {
// 	qt := make([]QuoteType, 0)
// 	k, err := q.Indexes.offsetIndexFromWord(w)
// 	if err != nil {
// 		return qt
// 	}
//
// 	for _, u := range q.Indexes[*k].Uuids {
// 		j, _ := q.OffsetQuoteFromUUID(u.String())
// 		quote := q.Data[*j]
// 		qt = append(qt, quote)
// 	}
//
// 	return qt
// }

// unparser writes a slice into database.
func unparser(db database.Storage, quotes Quotes) error {
	writeJSON, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return db.Write(writeJSON)
}
