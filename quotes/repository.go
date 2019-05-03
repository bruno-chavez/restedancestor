package quotes

import (
	"log"

	"github.com/bruno-chavez/restedancestor/database"
	uuid "github.com/satori/go.uuid"
)

func NewRepository(db database.Database) Repository {
	return Repository{db}
}

type Repository struct {
	db database.Database
}

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
		u2, err := uuid.FromString(u)
		if err != nil {
			log.Fatalf("Fail to parse : %s", u)
		}
		q := QuoteType{
			Quote: c,
			Uuid:  u2,
			Score: s,
		}

		quotes = append(quotes, q)
	}

	return quotes
}
