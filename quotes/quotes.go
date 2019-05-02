package quotes

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bruno-chavez/restedancestor/database"
	"github.com/satori/go.uuid"
)

// init is used to seed the rand.Intn function.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// index is an inverted index : for each word, list all documents that contain it.
type index struct {
	Word  string      `json:"word"`
	Uuids []uuid.UUID `json:"uuids"`
}

func (i index) isUUIDKnown(u uuid.UUID) bool {
	for _, uu := range i.Uuids {
		if uu == u {
			return true
		}
	}

	return false
}

type indexes []index

// offsetQuoteFromWord find the index with a word in the slice and returns its offset
func (i *indexes) offsetIndexFromWord(w string) (*int, error) {
	wLower := strings.ToLower(w)
	for k, idx := range *i {
		if idx.Word == wLower {
			return &k, nil
		}
	}

	return nil, errors.New("unknown index")
}

// setIndex creates or append an index
func (i *indexes) setIndex(w string, u uuid.UUID) {
	// First Find offset
	k, err := i.offsetIndexFromWord(w)

	// Doesn't exist, create index
	if err != nil {
		idx := index{
			Word:  w,
			Uuids: []uuid.UUID{u},
		}
		*i = append(*i, idx)
		return
	}

	// Exists, append u
	idx := &((*i)[*k])
	if !idx.isUUIDKnown(u) {
		idx.Uuids = append(idx.Uuids, u)
	}
}

// QuoteType describes a quote.
type QuoteType struct {
	Quote string    `json:"quote"`
	Uuid  uuid.UUID `json:"uuid"`
	Score int       `json:"score"`
}

// Quotes is used to parse the whole json database.
type Quotes struct {
	Data    []QuoteType `json:"data"`
	Indexes indexes     `json:"index"`
}

// Random returns a random quote from a Quotes type.
func (q Quotes) Random() QuoteType {
	qd := q.Data
	return qd[rand.Intn(len(qd))]
}

func (q Quotes) Len() int {
	return len(q.Data)
}

func (q Quotes) Swap(i int, j int) {
	qd := q.Data
	qd[i], qd[j] = qd[j], qd[i]
}

func (q Quotes) Less(i int, j int) bool {
	qd := q.Data
	return qd[i].Score > qd[j].Score
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

// LikeQuote increments the score of the quote
func (q Quotes) LikeQuote(db database.Storage, uuid string) {
	offset, _ := q.OffsetQuoteFromUUID(uuid)
	q.Data[*offset].Score++

	if err := unparser(db, q); err != nil {
		log.Fatal(err)
	}
}

// DislikeQuote decrements the score of the quote
func (q Quotes) DislikeQuote(db database.Storage, uuid string) {
	offset, _ := q.OffsetQuoteFromUUID(uuid)
	q.Data[*offset].Score--

	if err := unparser(db, q); err != nil {
		log.Fatal(err)
	}
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

// Index indexes all the data
func (q Quotes) Index(db database.Storage) {
	q.Indexes = make(indexes, 0)
	const limitSize = 3
	for _, quote := range q.Data {
		words := strings.FieldsFunc(quote.Quote, func(r rune) bool {
			switch r {
			case '\'', '!', ',', '.', ' ':
				return true
			}
			return false
		})
		for _, w := range words {
			wLower := strings.ToLower(w)
			if len(wLower) > limitSize {
				q.Indexes.setIndex(wLower, quote.Uuid)
			}
		}
	}

	unparser(db, q)
}

// List returns all quotes containing a word
func (q Quotes) List(w string) []QuoteType {
	qt := make([]QuoteType, 0)
	k, err := q.Indexes.offsetIndexFromWord(w)
	if err != nil {
		return qt
	}

	for _, u := range q.Indexes[*k].Uuids {
		j, _ := q.OffsetQuoteFromUUID(u.String())
		quote := q.Data[*j]
		qt = append(qt, quote)
	}

	return qt
}

// unparser writes a slice into database.
func unparser(db database.Storage, quotes Quotes) error {
	writeJSON, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return db.Write(writeJSON)
}
