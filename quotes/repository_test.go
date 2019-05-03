package quotes

import (
	// "errors"
	"testing"

	"github.com/bruno-chavez/restedancestor/database"
)

type DbDouble struct{}

type StmtDouble struct{}

func (d DbDouble) Prepare(sql string, args ...interface{}) (database.Stmt, error) {
	return StmtDouble{}, nil
}

func (s StmtDouble) Close() error {
	return nil
}

var step int

func (s StmtDouble) Step() (bool, error) {
	step++
	return (step <= 2), nil
}

var exec error

func (s StmtDouble) Exec(...interface{}) error {
	return exec
}

func (s StmtDouble) Scan(dst ...interface{}) error {
	return nil
}

var repo = NewRepository(DbDouble{})

func TestAllOK(t *testing.T) {
	step = 1
	ps := repo.All()
	if len(ps) != 1 {
		t.Error("No quote")
	}
}

func TestAllKO(t *testing.T) {
	step = 3
	ps := repo.All()
	if len(ps) != 0 {
		t.Error("There's quotes")
	}
}
