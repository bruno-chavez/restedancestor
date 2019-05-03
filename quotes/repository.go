package quotes

import (
	"errors"
	"log"

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

// All returns all quotes
func (r Repository) All() []QuoteType {
	stmt, err := r.db.Prepare(`SELECT content, score, uuid
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
	stmt, err := r.db.Prepare(`SELECT content, score, uuid
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

		var c string
		var s int
		var u string
		// q := QuoteType{}

		err = stmt.Scan(&c, &s, &u)
		if err != nil {
			log.Fatal("Scan gave error :" + err.Error())
		}

		// Parsing UUID from string input
		u2, _ := uuid.FromString(u)
		q := QuoteType{
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
