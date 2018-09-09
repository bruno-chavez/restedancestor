package database_test

import (
	"io/ioutil"
	"testing"

	d "github.com/bruno-chavez/restedancestor/database"
)

func TestRandom(t *testing.T) {
	s := make(d.QuoteSlice, 0)
	q := d.QuoteType{
		Quote: "this is a test",
	}
	s = append(s, q)
	rand := s.Random()
	if s.Random() != q {
		t.Errorf("Unexpected value %s, expected %s", rand, q)
	}
}

func TestNewDb(t *testing.T) {
	db := d.NewDb("/tmp/db")

	if len(db.Path()) == 0 {
		t.Error("Unexpected empty value for path")
	}
}

func TestParser(t *testing.T) {
	path := "/tmp/db"
	json := []byte(`[{"quote":"To be or not to be."}]`)
	if err := ioutil.WriteFile(path, json, 0644); err != nil {
		t.Fatalf("Unable to write into test file : %s", err.Error())
	}

	db := d.NewDb(path)
	if len(db.Parser()) == 0 {
		t.Error("Unexpected empty slice")
	}
}
