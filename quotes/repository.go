package quotes

import (
	"errors"
	"log"
	"strings"

	"github.com/bruno-chavez/restedancestor/database"
	uuid "github.com/satori/go.uuid"
)

// NewRepository initialises a new repo, for interacting with storage
func NewRepository(db database.Database) Repository {
	return Repository{db}
}

// Repository represents a storage abstraction
// Cf. https://www.martinfowler.com/eaaCatalog/repository.html
type Repository struct {
	db database.Database
}

// Random returns one quote, randomly
func (r Repository) Random() *QuoteType {
	stmt, err := r.db.Prepare(`SELECT id_quote, content, score, uuid
        FROM quotes
        ORDER BY RANDOM()
        LIMIT 1`)

	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	slice := buildSliceFromData(stmt)
	if len(slice) == 0 {
		return nil
	}

	return &slice[0]
}

// All returns all quotes
func (r Repository) All() []QuoteType {
	stmt, err := r.db.Prepare(`SELECT id_quote, content, score, uuid
        FROM quotes
        ORDER BY id_quote`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	return buildSliceFromData(stmt)
}

// FindByUUID returns one quote, given its UUID
func (r Repository) FindByUUID(u string) *QuoteType {
	stmt, err := r.Db.Prepare(`SELECT id_quote content, score, uuid
        FROM quotes
        WHERE uuid = ?`, u)

	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	slice := buildSliceFromData(stmt)
	if len(slice) == 0 {
		return nil
	}

	return &slice[0]
}

// Prefered returns 5 prefered quotes
func (r Repository) Prefered() []QuoteType {
	stmt, err := r.Db.Prepare(`SELECT id_quote, content, score, uuid
        FROM quotes
        ORDER BY score DESC
        LIMIT 5`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	return buildSliceFromData(stmt)
}

func buildSliceFromData(stmt database.Stmt) []QuoteType {
	quotes := make([]QuoteType, 0)

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			log.Fatal("Step gave error :" + err.Error())
		}
		if !hasRow {
			break
		}

		var i int
		var c string
		var s int
		var u string
		// q := QuoteType{}

		err = stmt.Scan(&i, &c, &s, &u)
		if err != nil {
			log.Fatal("Scan gave error :" + err.Error())
		}

		// Parsing UUID from string input
		u2, _ := uuid.FromString(u)
		q := QuoteType{
			ID:    i,
			Quote: c,
			Uuid:  u2,
			Score: s,
		}

		quotes = append(quotes, q)
	}

	return quotes
}

func (r Repository) IncrementsScore(u string) error {
	stmt, err := r.db.Prepare(`UPDATE quotes SET score = score+1
        WHERE uuid = ?`, u)
	if err != nil {
		return errors.New("Failed to prepare :" + err.Error())
	}
	defer stmt.Close()

	if err = stmt.Exec(u); err != nil {
		return errors.New("Failed to exec SQL :" + err.Error())
	}

	return nil
}

// @TODO : check existence, error message
func (r Repository) DecrementsScore(u string) error {
	stmt, err := r.db.Prepare(`UPDATE quotes SET score = score-1
        WHERE uuid = ?`, u)
	if err != nil {
		return errors.New("Failed to prepare :" + err.Error())
	}
	defer stmt.Close()

	if err = stmt.Exec(u); err != nil {
		return errors.New("Failed to exec SQL :" + err.Error())
	}

	return nil
}

func (r Repository) Index(q QuoteType) {
	const limitSize = 3

	words := strings.FieldsFunc(q.Quote, func(r rune) bool {
		switch r {
		case '\'', '!', ',', '.', '-', '…', ' ':
			return true
		}
		return false
	})
	for _, w := range words {
		wLower := strings.ToLower(w)
		if len(wLower) > limitSize {
			r.storeIndex(wLower, q.ID)
			log.Println("Store :", wLower, q.ID)
		}
	}
}

func (r Repository) storeIndex(w string, i int) {
	idWord := r.getIndexIDFromWord(w)
	stmt, err := r.Db.Prepare(`INSERT INTO indexes_quotes
        (id_index, id_quote) VALUES (?, ?)`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	err = stmt.Exec(idWord, i)

	if err != nil {
		log.Fatal("Failed to exec SQL :" + err.Error())
	}
}

func (r Repository) getIndexIDFromWord(w string) int64 {
	// log.Fatal(w)
	stmt, err := r.Db.Prepare(`SELECT id_index
        FROM indexes
        WHERE word = ?
        LIMIT 1`, w)

	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	hasRow, err := stmt.Step()
	if err != nil {
		log.Fatal("Step gave error :" + err.Error())
	}
	if !hasRow {
		return r.setIndexIDFromWord(w)
	}

	var idIndex int64
	err = stmt.Scan(&idIndex)
	if err != nil {
		log.Fatal("Scan gave error :" + err.Error())
	}
	return idIndex
}

func (r Repository) setIndexIDFromWord(w string) int64 {
	stmt, err := r.Db.Prepare(`INSERT INTO indexes
        (word) VALUES (?)`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	err = stmt.Exec(w)

	if err != nil {
		log.Fatal("Failed to exec SQL :" + err.Error())
	}

	// log.Fatal("Fatal :", r.Db.LastInsertRowID())

	return r.Db.LastInsertRowID()
}
// FindByWord()
// Index()
