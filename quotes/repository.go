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
func (r Repository) Random() *Quote {
	stmt, err := r.db.Prepare(`SELECT id_quote, content, score, uuid
        FROM quotes
        ORDER BY RANDOM()
        LIMIT 1`)

	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	slice := buildSliceFromData(stmt)
	if len(slice) == 0 {
		return nil
	}

	return &slice[0]
}

// All returns all quotes
func (r Repository) All() []Quote {
	stmt, err := r.db.Prepare(`SELECT id_quote, content, score, uuid
        FROM quotes
        ORDER BY id_quote`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	return buildSliceFromData(stmt)
}

// AllByWord returns all quotes containing a specific word
func (r Repository) AllByWord(w string) []Quote {
	stmt, err := r.db.Prepare(`SELECT DISTINCT q.id_quote, content, score, uuid
        FROM indexes i INNER JOIN indexes_quotes iq ON i.id_index = iq.id_index INNER JOIN quotes q ON iq.id_quote = q.id_quote
        WHERE word = ?
        ORDER BY q.id_quote`, w)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return buildSliceFromData(stmt)
}

// FindByUUID returns one quote, given its UUID
func (r Repository) FindByUUID(u string) *Quote {
	stmt, err := r.db.Prepare(`SELECT id_quote, content, score, uuid
        FROM quotes
        WHERE uuid = ?`, u)

	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	slice := buildSliceFromData(stmt)
	if len(slice) == 0 {
		return nil
	}

	return &slice[0]
}

// Preferred returns 5 preferred quotes
func (r Repository) Preferred() []Quote {
	stmt, err := r.db.Prepare(`SELECT id_quote, content, score, uuid
        FROM quotes
        ORDER BY score DESC
        LIMIT 5`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return buildSliceFromData(stmt)
}

func buildSliceFromData(stmt database.Stmt) []Quote {
	quotes := make([]Quote, 0)

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

		err = stmt.Scan(&i, &c, &s, &u)
		if err != nil {
			log.Fatal("Scan gave error :" + err.Error())
		}

		// Parsing UUID from string input
		u2, _ := uuid.FromString(u)
		q := Quote{
			id:    i,
			Quote: c,
			UUID:  u2,
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

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
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

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err = stmt.Exec(u); err != nil {
		return errors.New("Failed to exec SQL :" + err.Error())
	}

	return nil
}

// ----------
// Methods kept for indexation purpose

func (r Repository) index(q Quote) {
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
			r.storeIndex(wLower, q.id)
			log.Println("Store :", wLower, q.id)
		}
	}
}

func (r Repository) storeIndex(w string, i int) {
	idWord := r.getIndexIDFromWord(w)
	stmt, err := r.db.Prepare(`INSERT INTO indexes_quotes
        (id_index, id_quote) VALUES (?, ?)`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = stmt.Exec(idWord, i)

	if err != nil {
		log.Fatal("Failed to exec SQL :" + err.Error())
	}
}

func (r Repository) getIndexIDFromWord(w string) int64 {
	stmt, err := r.db.Prepare(`SELECT id_index
        FROM indexes
        WHERE word = ?
        LIMIT 1`, w)

	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

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
	stmt, err := r.db.Prepare(`INSERT INTO indexes
        (word) VALUES (?)`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}

	// closes db connection
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = stmt.Exec(w)

	if err != nil {
		log.Fatal("Failed to exec SQL :" + err.Error())
	}

	return r.db.LastInsertRowID()
}
